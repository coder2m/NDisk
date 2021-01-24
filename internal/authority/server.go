package authority

import (
	"context"
	"github.com/BurntSushi/toml"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xcode"
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/pkg/xjson"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xcfg/datasource/manager"
	"github.com/myxy99/component/xgovern"
	"github.com/myxy99/component/xinvoker"
	xgorm "github.com/myxy99/component/xinvoker/gorm"
	"github.com/myxy99/component/xmonitor"
	xclient "github.com/myxy99/ndisk/internal/authority/client"
	"github.com/myxy99/ndisk/internal/authority/model"
	auth_server "github.com/myxy99/ndisk/internal/authority/server"
	myValidator "github.com/myxy99/ndisk/internal/authority/validator"
	"github.com/myxy99/ndisk/pkg/rpc"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Server struct {
	err error
	*sync.WaitGroup
}

func (s *Server) PrepareRun(stopCh <-chan struct{}) (err error) {
	s.initCfg()
	s.invoker()
	s.debug()
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
		grpcCfg *xrpc.GRPCConfig
	)
	grpcCfg = xcfg.UnmarshalWithExpect("rpc", xrpc.DefaultGRPCConfig()).(*xrpc.GRPCConfig)
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
		xconsole.Red("grpc server shutdown success ")
		return nil
	})
	//NFilePb.RegisterNFileServiceServer(serve, new(rpcError.Server))
	xconsole.Greenf("grpc server start up success:", grpcCfg.Addr())
	s.err = serve.Serve(lis)
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
	xcode.GovernRun()
	xmonitor.Run()
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
	data, _ := auth_server.GetPermissionAndMenuByRoles(context.Background(), "3")
	st, _ := xjson.Marshal(data)
	xconsole.Red(string(st))
}
