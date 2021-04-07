package xclient

import (
	NUserPb "github.com/coder2z/ndisk/pkg/pb/nuser"
	xrpc "github.com/coder2z/ndisk/pkg/rpc"
)

var (
	userClient NUserPb.NUserServiceClient
)

func GetUserRpc() NUserPb.NUserServiceClient {
	if userClient == nil {
		userClient = NUserPb.NewNUserServiceClient(xrpc.Connection("ndisk_nuser"))
	}
	return userClient
}