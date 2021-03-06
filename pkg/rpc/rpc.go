package xrpc

import (
	"fmt"
	xapp "github.com/coder2z/component"
	"github.com/coder2z/component/xgrpc"
	clientinterceptors "github.com/coder2z/component/xgrpc/client"
	serverinterceptors "github.com/coder2z/component/xgrpc/server"
	"github.com/coder2z/component/xregistry"
	"github.com/coder2z/component/xregistry/xetcd"
	"github.com/coder2z/g-saber/xdefer"
	"github.com/coder2z/g-saber/xnet"
	"github.com/coder2z/ndisk/pkg/constant"
	"google.golang.org/grpc"
	"time"
)

type GRPCConfig struct {
	ServerIp      string        `mapStructure:"serverIp"`
	ServerPort    int           `mapStructure:"serverPort"`
	ServerName    string        `mapStructure:"serverName"`
	ServerTimeout time.Duration `mapStructure:"serverTimeout"`
	SlowThreshold time.Duration `mapStructure:"serverSlowThreshold"`

	EtcdAddr         string        `mapStructure:"register_etcd_addr"`
	RegisterTTL      time.Duration `mapStructure:"register_ttl"`
	RegisterInterval time.Duration `mapStructure:"register_interval"`
}

func DefaultGRPCConfig() *GRPCConfig {
	host, port, err := xnet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}
	return &GRPCConfig{
		ServerIp:         host,
		ServerPort:       port,
		ServerName:       xapp.Name(),
		ServerTimeout:    10 * time.Second,
		SlowThreshold:    8 * time.Second,
		EtcdAddr:         "127.0.0.1:2379",
		RegisterTTL:      30 * time.Second,
		RegisterInterval: 15 * time.Second,
	}
}

func (c GRPCConfig) Addr() string {
	return fmt.Sprintf("%v:%v", c.ServerIp, c.ServerPort)
}

func DefaultServerOption(c *GRPCConfig) []grpc.ServerOption {
	return []grpc.ServerOption{
		xgrpc.WithUnaryServerInterceptors(
			serverinterceptors.CrashUnaryServerInterceptor(),
			serverinterceptors.PrometheusUnaryServerInterceptor(),
			serverinterceptors.XTimeoutUnaryServerInterceptor(c.ServerTimeout),
			serverinterceptors.TraceUnaryServerInterceptor(),
		),
		xgrpc.WithStreamServerInterceptors(
			serverinterceptors.CrashStreamServerInterceptor(),
			serverinterceptors.PrometheusStreamServerInterceptor(),
		),
	}
}

func DefaultClientOption(c *GRPCConfig) []grpc.DialOption {
	return []grpc.DialOption{
		xgrpc.WithUnaryClientInterceptors(
			clientinterceptors.XTimeoutUnaryClientInterceptor(c.ServerTimeout, c.SlowThreshold),
			clientinterceptors.XTraceUnaryClientInterceptor(),
			clientinterceptors.XAidUnaryClientInterceptor(),
			clientinterceptors.XLoggerUnaryClientInterceptor(c.ServerName),
			clientinterceptors.PrometheusUnaryClientInterceptor(c.ServerName),
		),
		xgrpc.WithStreamClientInterceptors(
			clientinterceptors.PrometheusStreamClientInterceptor(c.ServerName),
		),
	}
}

func DefaultRegistryEtcd(c *GRPCConfig) (err error) {
	var etcdR xregistry.Registry
	conf := xetcd.EtcdV3Cfg{
		Endpoints: []string{c.EtcdAddr},
	}
	etcdR, err = xetcd.NewRegistry(conf) //注册
	if err != nil {
		return
	}

	etcdR.Register(
		xregistry.ServiceName(xapp.Name()),
		xregistry.ServiceNamespaces(constant.DefaultNamespaces),
		xregistry.Address(c.Addr()),
		xregistry.RegisterTTL(c.RegisterTTL),
		xregistry.RegisterInterval(c.RegisterInterval),
	)

	xdefer.Register(func() error {
		etcdR.Close()
		return nil
	})
	return
}
