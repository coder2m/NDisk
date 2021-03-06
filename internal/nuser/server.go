package nuser

import (
	"net"
	"sync"

	"github.com/BurntSushi/toml"
	xapp "github.com/coder2m/component"
	"github.com/coder2m/component/xcfg"
	"github.com/coder2m/component/xcfg/datasource/manager"
	"github.com/coder2m/component/xgovern"
	"github.com/coder2m/component/xinvoker"
	xemail "github.com/coder2m/component/xinvoker/email"
	xgorm "github.com/coder2m/component/xinvoker/gorm"
	xredis "github.com/coder2m/component/xinvoker/redis"
	xsms "github.com/coder2m/component/xinvoker/sms"
	"github.com/coder2m/component/xmonitor"
	"github.com/coder2m/component/xtrace"
	"github.com/coder2m/g-saber/xconsole"
	"github.com/coder2m/g-saber/xdefer"
	"github.com/coder2m/g-saber/xflag"
	"github.com/coder2m/g-saber/xvalidator"
	"github.com/coder2m/ndisk/internal/nuser/model"
	"github.com/coder2m/ndisk/internal/nuser/rpc"
	myValidator "github.com/coder2m/ndisk/internal/nuser/validator"
	NUserPb "github.com/coder2m/ndisk/pkg/pb/nuser"
	"github.com/coder2m/ndisk/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	xtrace.Init("trace.jaeger")
	go xgovern.Run()
}
