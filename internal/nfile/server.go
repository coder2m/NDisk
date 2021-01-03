package nfile

import (
	"context"
	"github.com/BurntSushi/toml"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xcfg/datasource/manager"
	"github.com/myxy99/component/xgovern"
	"github.com/myxy99/component/xgrpc"
	serverinterceptors "github.com/myxy99/component/xgrpc/server"
	"github.com/myxy99/component/xinvoker"
	xgorm "github.com/myxy99/component/xinvoker/gorm"
	"github.com/myxy99/component/xmonitor"
	"github.com/myxy99/component/xregistry"
	"github.com/myxy99/component/xregistry/xetcd"
	"github.com/myxy99/ndisk/internal/nfile/api/v1/registry"
	"github.com/myxy99/ndisk/internal/nfile/rpc"
	myValidator "github.com/myxy99/ndisk/internal/nfile/validator"
	"github.com/myxy99/ndisk/pkg/constant"
	NFilePb "github.com/myxy99/ndisk/pkg/pb/nfile"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	Server *http.Server
	err    error
	*sync.WaitGroup
}

func (s *Server) PrepareRun(stopCh <-chan struct{}) (err error) {
	s.initCfg()
	s.debug()
	s.initHttpServer()
	s.initRouter()
	s.rpc()
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
	xdefer.Register(func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		xconsole.Red("http server shutdown")
		return s.Server.Shutdown(ctx)
	})
	xconsole.Greenf("Start listening on:", s.Server.Addr)
	if err = s.Server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	s.Wait()
	return nil
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
	)
	s.err = xinvoker.Init()
}

func (s *Server) initHttpServer() {
	if s.err != nil {
		return
	}
	s.Server = new(http.Server)
	s.Server.Addr = xcfg.GetString("server.addr")
}

func (s *Server) initRouter() {
	if s.err != nil {
		return
	}
	s.Server.Handler = registry.Engine()
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
	go xgovern.Run()
}

func (s *Server) rpc() {
	if s.err != nil {
		return
	}
	var (
		etcdR  xregistry.Registry
		rpcCfg *rpc.Config
		lis    net.Listener
	)
	rpcCfg = xcfg.UnmarshalWithExpect("rpc", rpc.DefaultConfig()).(*rpc.Config)
	conf := xetcd.EtcdV3Cfg{
		Endpoints: []string{rpcCfg.EtcdAddr},
	}
	etcdR, s.err = xetcd.NewRegistry(conf) //注册
	if s.err != nil {
		return
	}

	etcdR.Register(
		xregistry.ServiceName(xapp.Name()),
		xregistry.ServiceNamespaces(constant.DefaultNamespaces),
		xregistry.Address(rpcCfg.Addr()),
		xregistry.RegisterTTL(rpcCfg.RegisterTTL),
		xregistry.RegisterInterval(rpcCfg.RegisterInterval),
	)

	xdefer.Register(func() error {
		etcdR.Close()
		return nil
	})

	lis, s.err = net.Listen("tcp", rpcCfg.Addr())
	if s.err != nil {
		return
	}

	options := []grpc.ServerOption{
		xgrpc.WithUnaryServerInterceptors(
			serverinterceptors.CrashUnaryServerInterceptor(),
			serverinterceptors.PrometheusUnaryServerInterceptor(),
			serverinterceptors.XTimeoutUnaryServerInterceptor(rpcCfg.Timeout),
			serverinterceptors.TraceUnaryServerInterceptor(),
		),
		xgrpc.WithStreamServerInterceptors(
			serverinterceptors.CrashStreamServerInterceptor(),
			serverinterceptors.PrometheusStreamServerInterceptor(),
		),
	}

	serve := grpc.NewServer(options...)
	NFilePb.RegisterNFileServiceServer(serve, new(rpc.Server))
	go func() {
		s.err = serve.Serve(lis)
	}()
	xdefer.Register(func() error {
		serve.Stop()
		xconsole.Red("grpc server shutdown success ")
		return nil
	})
	xconsole.Greenf("grpc server start up success:", rpcCfg.Addr())
}
