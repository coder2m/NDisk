package authority

import (
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xtrace"
	"net"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xconsole"
	"github.com/coder2z/g-saber/xdefer"
	"github.com/coder2z/g-saber/xflag"
	"github.com/coder2z/g-saber/xvalidator"
	"github.com/coder2z/g-server/datasource/manager"
	"github.com/coder2z/g-server/xapp"
	"github.com/coder2z/g-server/xgovern"
	"github.com/coder2z/g-server/xinvoker"
	xgorm "github.com/coder2z/g-server/xinvoker/gorm"
	"google.golang.org/grpc"

	xclient "github.com/coder2z/ndisk/internal/authority/client"
	"github.com/coder2z/ndisk/internal/authority/model"
	"github.com/coder2z/ndisk/internal/authority/rpc"
	myValidator "github.com/coder2z/ndisk/internal/authority/validator"
	AuthorityPb "github.com/coder2z/ndisk/pkg/pb/authority"
	"github.com/coder2z/ndisk/pkg/rpc"
)

type Server struct {
	err error
	*sync.WaitGroup
}

func (s *Server) PrepareRun(stopCh <-chan struct{}) (err error) {
	s.initCfg()
	s.debug()
	s.invoker()
	s.initValidator()
	s.govern()
	s.casbin()
	return s.err
}

func (s *Server) debug() {
	xconsole.ResetDebug(xapp.Debug())
	xapp.PrintVersion()
}

func (s *Server) Run(stopCh <-chan struct{}) (err error) {
	go func() {
		<-stopCh
		s.Add(1)
		xdefer.Clean()
		s.Done()
	}()
	var (
		lis     net.Listener
		grpcCfg *xrpc.GRPCServerConfig
	)
	grpcCfg = xrpc.GRPCServerCfgBuild("rpc")
	s.err = xrpc.DefaultRegistryEtcd(grpcCfg)
	if s.err != nil {
		return
	}
	lis, s.err = net.Listen("tcp", grpcCfg.Addr())
	if s.err != nil {
		return
	}
	serve := grpc.NewServer(xrpc.DefaultServerOption(grpcCfg)...)
	xdefer.Register(func() error {
		serve.Stop()
		xlog.Info("Application Stopping",
			xlog.FieldComponentName("GRPC"),
			xlog.FieldMethod("Auth.Run"),
			xlog.FieldDescription("GRPC server shutdown success"),
		)
		return nil
	})
	AuthorityPb.RegisterAuthorityServiceServer(serve, new(rpc.Server))
	s.err = serve.Serve(lis)
	xlog.Info("Application Starting",
		xlog.FieldComponentName("GRPC"),
		xlog.FieldMethod("Auth.Run"),
		xlog.FieldDescription("GRPC server start up success"),
		xlog.FieldAddr(grpcCfg.Addr()),
	)
	s.Wait()
	return s.err
}

func (s *Server) initCfg() {
	if s.err != nil {
		return
	}
	var data xcfg.DataSource
	data, s.err = manager.NewDataSource(xflag.NString("run", "xcfg"))
	if s.err != nil {
		return
	}
	s.err = xcfg.LoadFromDataSource(data, toml.Unmarshal)
}

func (s *Server) invoker() {
	if s.err != nil {
		return
	}
	xdefer.Register(func() error {
		return xinvoker.Close()
	})
	xinvoker.Register(
		xgorm.Register("mysql"),
		//xredis.Register("redis"),
	)
	s.err = xinvoker.Init()
	_ = model.MainDB().Migrator().AutoMigrate(
		new(model.Resources),
		new(model.Menu),
		new(model.Roles),
	)
}

func (s *Server) initValidator() {
	if s.err != nil {
		return
	}
	s.err = xvalidator.Init(xcfg.GetString("server.locale"), myValidator.RegisterValidation)
}

func (s *Server) govern() {
	if s.err != nil {
		return
	}
	xtrace.JaegerBuild("trace.jaeger")
	go xgovern.Run()
}

func (s *Server) casbin() {
	if s.err != nil {
		return
	}
	xclient.CasbinClient()
	// TODO 测试
	//xclient.CasbinClient().AddPermissionForUser("测试角色", "/test", "GET")
	//_ = auth_server.AddRoles(context.Background(), _map.RolesReq{
	//	Name:        "admin2",
	//	Description: "我是admin2",
	//})
	//st, _ := xjson.Marshal(data)
	//xconsole.Red(string(st))
	//
	//_ = auth_server.UpdateRolesMenuAndResources(context.Background(), _map.UpdateRolesMenuAndResourcesReq{
	//	ID:        3,
	//	Menus:     []uint32{1, 2},
	//	Resources: []uint32{1, 2},
	//})
	//
	//data, _ := auth_server.GetPermissionAndMenuByRoles(context.Background(), "3")
	//st, _ := xjson.Marshal(data)
	//xconsole.Red(string(st))
}
