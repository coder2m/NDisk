package xclient

import (
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/ndisk/pkg/constant"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
	"google.golang.org/grpc"
)

var (
	NUserServer NUserPb.NUserServiceClient
	grpcCfg     *xrpc.GRPCConfig
)

func InitClient() {
	NUserServer = NUserPb.NewNUserServiceClient(connection("ndisk_nuser"))
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
