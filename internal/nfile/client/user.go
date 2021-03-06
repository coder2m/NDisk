package xclient

import (
	"context"

	NUserPb "github.com/coder2z/ndisk/pkg/pb/nuser"
)

func GetUserInfoByToken(ctx context.Context, token string) (*NUserPb.UserInfo, error) {
	return GetUserRpc().VerifyUsers(context.Background(), &NUserPb.Token{
		AccountToken: token,
	})
}
