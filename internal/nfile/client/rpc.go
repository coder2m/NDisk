package xclient

import (
	"github.com/coder2m/component/pkg/xconsole"
	"github.com/coder2m/component/pkg/xdefer"
	"github.com/coder2m/component/xcfg"
	"google.golang.org/grpc"

	"github.com/coder2m/ndisk/pkg/constant"
	NUserPb "github.com/coder2m/ndisk/pkg/pb/nuser"
	xrpc "github.com/coder2m/ndisk/pkg/rpc"
)

var (
	userClient NUserPb.NUserServiceClient
	grpcCfg    *xrpc.GRPCConfig
)

func GetUserRpc() NUserPb.NUserServiceClient {
	if userClient == nil {
		userClient = NUserPb.NewNUserServiceClient(connection("ndisk_nuser"))
	}
	return userClient
}

func GetGRPCCfg() *xrpc.GRPCConfig {
	if grpcCfg == nil {
		grpcCfg = xcfg.UnmarshalWithExpect("rpc", xrpc.DefaultGRPCConfig()).(*xrpc.GRPCConfig)
	}
	return grpcCfg
}

func connection(servername string, op ...grpc.DialOption) *grpc.ClientConn {
	option := xrpc.DefaultClientOption(GetGRPCCfg())
	option = append(option, grpc.WithInsecure())
	option = append(option, op...)
	conn, err := grpc.Dial(constant.GRPCTargetEtcd.Format(constant.DefaultNamespaces, servername), option...)
	if err != nil {
		panic(err.Error())
	}
	xdefer.Register(func() error {
		xconsole.Redf("grpc conn close => server name:", servername)
		return conn.Close()
	})
	return conn
}
