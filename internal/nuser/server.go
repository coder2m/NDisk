package nuser

import (
	"github.com/BurntSushi/toml"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xcode"
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xcfg/datasource/manager"
	"github.com/myxy99/component/xgovern"
	"github.com/myxy99/component/xinvoker"
	xemail "github.com/myxy99/component/xinvoker/email"
	xgorm "github.com/myxy99/component/xinvoker/gorm"
	xredis "github.com/myxy99/component/xinvoker/redis"
	xsms "github.com/myxy99/component/xinvoker/sms"
	"github.com/myxy99/component/xmonitor"
	"github.com/myxy99/component/xtrace"
	"github.com/myxy99/ndisk/internal/nuser/model"
	"github.com/myxy99/ndisk/internal/nuser/rpc"
	myValidator "github.com/myxy99/ndisk/internal/nuser/validator"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	"github.com/myxy99/ndisk/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
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
	return s.err
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
	NUserPb.RegisterNUserServiceServer(serve, new(rpc.Server))
	xconsole.Greenf("grpc server start up success:", grpcCfg.Addr())
	if xapp.Debug() {
		reflection.Register(serve)
	}
	s.err = serve.Serve(lis)
	s.Wait()
	return s.err
}

func (s *Server) debug() {
	xconsole.ResetDebug(xapp.Debug())
	xapp.PrintVersion()
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
		xredis.Register("redis"),
		xemail.Register("email"),
		xsms.Register("sms"),
	)
	s.err = xinvoker.Init()
	_ = model.MainDB().Migrator().AutoMigrate(
		new(model.User),
		new(model.Agency),
		new(model.AgencyUser),
	)
	_ = model.MainDB().SetupJoinTable(&model.User{}, "Agency", &model.AgencyUser{})
	_ = model.MainDB().SetupJoinTable(&model.Agency{}, "Users", &model.AgencyUser{})
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
	xmonitor.Run()
	xcode.GovernRun()
	xtrace.Init("trace.jaeger")
	go xgovern.Run()
}
