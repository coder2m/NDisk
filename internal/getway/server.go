package getway

import (
	"context"
	"net/http"
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
	"github.com/coder2z/g-server/xregistry/xetcd"
	"github.com/coder2z/g-server/xtrace"
	"github.com/coder2z/ndisk/internal/getway/api/v1/registry"
	"github.com/coder2z/ndisk/internal/getway/client"
	myValidator "github.com/coder2z/ndisk/internal/getway/validator"
	"github.com/coder2z/ndisk/pkg/recaptcha"
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
	s.initValidator()
	s.govern()
	s.rpc()
	s.recaptcha()
	return s.err
}

func (s *Server) Run(stopCh <-chan struct{}) (err error) {
	go func() {
		defer s.Done()
		<-stopCh
		s.Add(1)
		xdefer.Clean()
	}()
	xdefer.Register(func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		xconsole.Red("http server Shutdown")
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
	xtrace.JaegerBuild("trace.jaeger")
	go xgovern.Run()
}

func (s *Server) rpc() {
	if s.err != nil {
		return
	}
	grpcCfg := xclient.GetGRPCCfg()
	conf := xetcd.EtcdV3Cfg{
		Endpoints:        []string{grpcCfg.EtcdAddr},
		AutoSyncInterval: grpcCfg.RegisterInterval,
	}
	s.err = xetcd.RegisterBuilder(conf)
	xclient.InitClient()
}

func (s *Server) recaptcha() {
	recaptchaCfg := xcfg.UnmarshalWithExpect("google.recaptcha", recaptcha.DefaultConfig()).(*recaptcha.Config)
	recaptcha.NewDefault(recaptchaCfg)
}
