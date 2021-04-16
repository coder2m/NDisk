package xrpc

import (
	"fmt"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xdefer"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xnet"
	"github.com/coder2z/g-server/xapp"
	"github.com/coder2z/g-server/xgrpc"
	xbalancer "github.com/coder2z/g-server/xgrpc/balancer"
	"github.com/coder2z/g-server/xgrpc/balancer/p2c"
	clientinterceptors "github.com/coder2z/g-server/xgrpc/client"
	serverinterceptors "github.com/coder2z/g-server/xgrpc/server"
	"github.com/coder2z/g-server/xregistry"
	"github.com/coder2z/g-server/xregistry/xetcd"
	"github.com/coder2z/ndisk/pkg/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

type GRPCServerConfig struct {
	ServerTimeout time.Duration `mapStructure:"serverTimeout"`
	ServerIp      string        `mapStructure:"serverIp"`
	ServerPort    int           `mapStructure:"serverPort"`
	ETCD          EtcdConfig    `mapStructure:"etcd"`
	Weight        string        `mapStructure:"weight"`
}

type GRPCClientConfig struct {
	ServerName    string        `mapStructure:"serverName"`
	ServerTimeout time.Duration `mapStructure:"serverTimeout"`
	SlowThreshold time.Duration `mapStructure:"serverSlowThreshold"`
	ETCD          EtcdConfig    `mapStructure:"etcd"`
}

type EtcdConfig struct {
	EtcdAddr         []string      `mapStructure:"register_etcd_addr"`
	RegisterTTL      time.Duration `mapStructure:"register_ttl"`
	RegisterInterval time.Duration `mapStructure:"register_interval"`
}

func DefaultGRPCClientConfig() *GRPCClientConfig {
	return &GRPCClientConfig{
		ServerName:    xapp.Name(),
		ServerTimeout: 10 * time.Second,
		SlowThreshold: 8 * time.Second,
		ETCD: EtcdConfig{
			EtcdAddr:         []string{"127.0.0.1:2379"},
			RegisterTTL:      30 * time.Second,
			RegisterInterval: 15 * time.Second,
		},
	}
}

func GRPCClientCfgBuild(key string) *GRPCClientConfig {
	grpcCfg := xcfg.UnmarshalWithExpect(key, DefaultGRPCClientConfig()).(*GRPCClientConfig)
	return grpcCfg
}

func GRPCServerCfgBuild(key string) *GRPCServerConfig {
	grpcCfg := xcfg.UnmarshalWithExpect(key, DefaultGRPCServerConfig()).(*GRPCServerConfig)
	return grpcCfg
}

func DefaultGRPCServerConfig() *GRPCServerConfig {
	host, port, err := xnet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}
	return &GRPCServerConfig{
		ServerTimeout: 10 * time.Second,
		ServerIp:      host,
		ServerPort:    port,
		ETCD: EtcdConfig{
			EtcdAddr:         []string{"127.0.0.1:2379"},
			RegisterTTL:      30 * time.Second,
			RegisterInterval: 15 * time.Second,
		},
		Weight: "1",
	}
}

func (c GRPCServerConfig) Addr() string {
	return fmt.Sprintf("%v:%v", c.ServerIp, c.ServerPort)
}

func DefaultServerOption(c *GRPCServerConfig) []grpc.ServerOption {
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

func DefaultClientOption(c *GRPCClientConfig) []grpc.DialOption {
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
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, p2c.P2C)),
	}
}

func DefaultRegistryEtcd(c *GRPCServerConfig) (err error) {
	var etcdR xregistry.Registry
	conf := xetcd.EtcdV3Cfg{
		Endpoints: c.ETCD.EtcdAddr,
	}
	etcdR, err = xetcd.NewRegistry(conf) //注册
	if err != nil {
		return
	}

	etcdR.Register(
		xregistry.ServiceName(xapp.Name()),
		xregistry.ServiceNamespaces(constant.DefaultNamespaces),
		xregistry.Address(c.Addr()),
		xregistry.RegisterTTL(c.ETCD.RegisterTTL),
		xregistry.RegisterInterval(c.ETCD.RegisterInterval),
		xregistry.Metadata(metadata.Pairs(xbalancer.WeightKey, c.Weight)),
	)

	err = RegistryBuilder(c.ETCD)

	xdefer.Register(func() error {
		etcdR.Close()
		return nil
	})
	return
}

func Connection(servername string, op ...grpc.DialOption) *grpc.ClientConn {
	cfg := GRPCClientCfgBuild("rpc." + servername)
	option := DefaultClientOption(cfg)
	option = append(option, grpc.WithInsecure())
	option = append(option, op...)
	conn, err := grpc.Dial(constant.GRPCTargetEtcd.Format(constant.DefaultNamespaces, servername), option...)
	if err != nil {
		panic(err.Error())
	}
	xdefer.Register(func() error {
		xlog.Info("Application Stopping",
			xlog.FieldComponentName("GRPC"),
			xlog.FieldMethod("RPC.Connection"),
			xlog.FieldDescription("GRPC connection close"),
			xlog.String("server name", servername),
		)
		return conn.Close()
	})
	return conn
}

func RegistryBuilder(etcd EtcdConfig) error {
	conf := xetcd.EtcdV3Cfg{
		Endpoints:        etcd.EtcdAddr,
		AutoSyncInterval: etcd.RegisterInterval,
	}
	return xetcd.RegisterBuilder(conf)
}
