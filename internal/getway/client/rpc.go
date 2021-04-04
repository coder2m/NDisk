package xclient

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xconsole"
	"github.com/coder2z/g-saber/xdefer"
	"github.com/coder2z/ndisk/pkg/constant"
	AuthorityPb "github.com/coder2z/ndisk/pkg/pb/authority"
	NUserPb "github.com/coder2z/ndisk/pkg/pb/nuser"
	xrpc "github.com/coder2z/ndisk/pkg/rpc"
	"google.golang.org/grpc"
)

var (
	NUserServer     NUserPb.NUserServiceClient
	AuthorityServer AuthorityPb.AuthorityServiceClient
	grpcCfg         *xrpc.GRPCConfig
)

func InitClient() {
	NUserServer = NUserPb.NewNUserServiceClient(connection("ndisk_nuser"))
	AuthorityServer = AuthorityPb.NewAuthorityServiceClient(connection("ndisk_authority"))
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
