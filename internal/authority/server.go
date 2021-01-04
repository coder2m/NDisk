package authority

import (
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
	"github.com/myxy99/ndisk/internal/authority/model"
	myValidator "github.com/myxy99/ndisk/internal/authority/validator"
	"github.com/myxy99/ndisk/pkg/constant"
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
	s.initDB(stopCh)
	s.initValidator()
	s.govern()
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
	//NFilePb.RegisterNFileServiceServer(serve, new(rpc.Server))
	s.err = serve.Serve(lis)
	if s.err != nil {
		return s.err
	}
	xdefer.Register(func() error {
		serve.Stop()
		xconsole.Red("grpc server shutdown success ")
		return nil
	})
	xconsole.Greenf("grpc server start up success:", rpcCfg.Addr())
	s.Wait()
	return nil
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
		//xgorm.Register("mysql"),
		//xredis.Register("redis"),
	)
	s.err = xinvoker.Init()
}

func (s *Server) initDB(stopCh <-chan struct{}) {
	if s.err != nil {
		return
	}
	model.MainDB = xgorm.Invoker("main")
	go func() {
		<-stopCh
		d, _ := model.MainDB.DB()
		_ = d.Close()
	}()
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
