package nfile

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xcfg/datasource/manager"
	"github.com/myxy99/component/xgovern"
	"github.com/myxy99/component/xinvoker"
	xgorm "github.com/myxy99/component/xinvoker/gorm"
	"ndisk/internal/nfile/api/v1/registry"
	"ndisk/internal/nfile/model"
	myValidator "ndisk/internal/nfile/validator"
	"net/http"
)

type Server struct {
	Server *http.Server
	err    error
}

func (s *Server) PrepareRun(stopCh <-chan struct{}) (err error) {
	s.initCfg()
	s.invoker()
	s.initDB(stopCh)
	s.initHttpServer()
	s.initRouter()
	s.initValidator()
	s.govern()
	return s.err
}

func (s *Server) Run(stopCh <-chan struct{}) (err error) {
	go func() {
		<-stopCh
		xdefer.Clean()
	}()
	xdefer.Register(func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		xconsole.Red("http server Shutdown")
		return s.Server.Shutdown(ctx)
	})
	xconsole.Greenf("Start listening on:", s.Server.Addr)
	err = s.Server.ListenAndServe()
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
		xgorm.Register("mysql"),
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

	s.initMigrate()
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

func (s *Server) initMigrate() {
	if s.err != nil {
		return
	}
	if model.MainDB != nil {
		_ = model.MainDB.AutoMigrate(
			//new(info.Info),
		)
	}
}

func (s *Server) govern() {
	if s.err != nil {
		return
	}
	go xgovern.Run()
}
