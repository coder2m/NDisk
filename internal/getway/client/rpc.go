package xclient

import (
	AuthorityPb "github.com/coder2z/ndisk/pkg/pb/authority"
	NUserPb "github.com/coder2z/ndisk/pkg/pb/nuser"
	xrpc "github.com/coder2z/ndisk/pkg/rpc"
)

var (
	NUserServer     NUserPb.NUserServiceClient
	AuthorityServer AuthorityPb.AuthorityServiceClient
)

func InitClient() {
	NUserServer = NUserPb.NewNUserServiceClient(xrpc.Connection("ndisk_nuser"))
	AuthorityServer = AuthorityPb.NewAuthorityServiceClient(xrpc.Connection("ndisk_auth"))
}