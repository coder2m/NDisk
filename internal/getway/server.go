package getway

import (
	"context"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xinvoker"
	xgorm "github.com/coder2z/g-server/xinvoker/gorm"
	xredis "github.com/coder2z/g-server/xinvoker/redis"
	"github.com/coder2z/ndisk/internal/getway/model"
	xrpc "github.com/coder2z/ndisk/pkg/rpc"
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
	s.invoker()
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
	grpcCfg := xrpc.GRPCServerCfgBuild("rpc")
	s.err = xrpc.RegistryBuilder(grpcCfg.ETCD)
	xclient.InitClient()
}

func (s *Server) recaptcha() {
	recaptchaCfg := xcfg.UnmarshalWithExpect("google.recaptcha", recaptcha.DefaultConfig()).(*recaptcha.Config)
	recaptcha.NewDefault(recaptchaCfg)
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
	)
	s.err = xinvoker.Init()
	db := model.MainDB()
	rdb := model.MainRedis()
	xlog.Infow("DB AutoMigrate", xlog.FieldErr(db.AutoMigrate(new(model.Directory))))
	xlog.Info("Redis ping", xlog.FieldErr(rdb.Ping(context.Background()).Err()))
}
