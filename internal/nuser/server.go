package nuser

import (
	"github.com/coder2z/g-saber/xlog"
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
	xemail "github.com/coder2z/g-server/xinvoker/email"
	xgorm "github.com/coder2z/g-server/xinvoker/gorm"
	xredis "github.com/coder2z/g-server/xinvoker/redis"
	xsms "github.com/coder2z/g-server/xinvoker/sms"
	"github.com/coder2z/g-server/xtrace"
	"github.com/coder2z/ndisk/internal/nuser/model"
	"github.com/coder2z/ndisk/internal/nuser/rpc"
	myValidator "github.com/coder2z/ndisk/internal/nuser/validator"
	NUserPb "github.com/coder2z/ndisk/pkg/pb/nuser"
	"github.com/coder2z/ndisk/pkg/rpc"
	"google.golang.org/grpc"
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
			xlog.FieldMethod("NUser.Run"),
			xlog.FieldDescription("GRPC server shutdown success"),
		)
		return nil
	})
	NUserPb.RegisterNUserServiceServer(serve, new(rpc.Server))
	s.err = serve.Serve(lis)
	xlog.Info("Application Starting",
		xlog.FieldComponentName("GRPC"),
		xlog.FieldMethod("NUser.Run"),
		xlog.FieldDescription("GRPC server start up success"),
		xlog.FieldAddr(grpcCfg.Addr()),
	)
	if s.err != nil {
		return
	}
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
	xtrace.JaegerBuild("trace.jaeger")
	go xgovern.Run()
}
