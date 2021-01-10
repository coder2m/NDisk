package getway

import (
	"context"
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
	"github.com/myxy99/component/xmonitor"
	"github.com/myxy99/ndisk/internal/getway/api/v1/registry"
	myValidator "github.com/myxy99/ndisk/internal/getway/validator"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
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
	s.initValidator()
	s.govern()
	s.rpc()
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
	xcode.GovernRun()
	xmonitor.Run()
	go xgovern.Run()
}

func (s *Server) rpc() {
	if s.err != nil {
		return
	}
	var rpcCfg *xrpc.GRPCConfig
	rpcCfg = xcfg.UnmarshalWithExpect("rpc", xrpc.DefaultGRPCConfig()).(*xrpc.GRPCConfig)
	s.err = xrpc.DefaultRegistryEtcd(rpcCfg)
}
