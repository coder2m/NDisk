package nfile

import (
	"context"
	xclient "github.com/coder2z/ndisk/internal/nfile/client"
	"net"
	"net/http"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xconsole"
	"github.com/coder2z/g-saber/xdefer"
	"github.com/coder2z/g-saber/xflag"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xvalidator"
	"github.com/coder2z/g-server/datasource/manager"
	"github.com/coder2z/g-server/xapp"
	"github.com/coder2z/g-server/xgovern"
	"github.com/coder2z/g-server/xinvoker"
	xgorm "github.com/coder2z/g-server/xinvoker/gorm"
	xoss "github.com/coder2z/g-server/xinvoker/oss"
	xredis "github.com/coder2z/g-server/xinvoker/redis"
	"google.golang.org/grpc"

	"github.com/coder2z/ndisk/internal/nfile/api/v1/registry"
	"github.com/coder2z/ndisk/internal/nfile/model"
	rpcServer "github.com/coder2z/ndisk/internal/nfile/rpc"
	myValidator "github.com/coder2z/ndisk/internal/nfile/validator"
	NFilePb "github.com/coder2z/ndisk/pkg/pb/nfile"
	"github.com/coder2z/ndisk/pkg/rpc"
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
		xredis.Register("redis"),
		xoss.Register("oss"),
	)
	s.err = xinvoker.Init()
	db := model.MainDB()
	xlog.Infow("AutoMigrate", "model.File", db.AutoMigrate(&model.File{}, &model.FileSlice{}))
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

	go xgovern.Run()
}

func (s *Server) rpc() {
	if s.err != nil {
		return
	}
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
	NFilePb.RegisterNFileServiceServer(serve, new(rpcServer.Server))
	go func() {
		s.err = serve.Serve(lis)
	}()
	xdefer.Register(func() error {
		serve.Stop()
		xconsole.Red("grpc server shutdown success ")
		return nil
	})
	xclient.GetUserRpc()
	xconsole.Greenf("grpc server start up success:", grpcCfg.Addr())
}
