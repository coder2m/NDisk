package client

import (
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xregistry/xetcd"
	"github.com/myxy99/ndisk/pkg/constant"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
	"google.golang.org/grpc"
)

var (
	nUserServer NUserPb.NUserServiceClient
	grpcCfg     *xrpc.GRPCConfig
)

func NUserServer() NUserPb.NUserServiceClient {
	if nUserServer == nil {
		option := xrpc.DefaultClientOption(GetGRPCCfg())
		option = append(option, grpc.WithInsecure())
		conn, err := grpc.Dial(constant.GRPCTargetEtcd.Format(constant.DefaultNamespaces, "ndisk_nuser"), option...)
		if err != nil {
			panic(err.Error())
		}

		xdefer.Register(func() error {
			xconsole.Red("grpc conn close")
			return conn.Close()
		})

		nUserServer = NUserPb.NewNUserServiceClient(conn)
	}
	return nUserServer
}

func GetGRPCCfg() *xrpc.GRPCConfig {
	if grpcCfg == nil {
		grpcCfg = xcfg.UnmarshalWithExpect("rpc", xrpc.DefaultGRPCConfig()).(*xrpc.GRPCConfig)
		conf := xetcd.EtcdV3Cfg{
			Endpoints:        []string{grpcCfg.EtcdAddr},
			AutoSyncInterval: grpcCfg.RegisterInterval,
		}
		_ = xetcd.RegisterBuilder(conf)
	}
	return grpcCfg
}
