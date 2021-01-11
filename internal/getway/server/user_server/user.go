/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 19:21
 **/
package user_server

import (
	"context"
	"errors"
	xclient "github.com/myxy99/ndisk/internal/getway/client"
	xerror "github.com/myxy99/ndisk/internal/getway/error"
	_map "github.com/myxy99/ndisk/internal/getway/map"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
)

func AccountLogin(ctx context.Context, login _map.AccountLogin) *xerror.Err {
	rep, err := xclient.NUserServer.AccountLogin(ctx, &NUserPb.UserLoginRequest{
		Account:  login.Account,
		Password: login.Password,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}

}
