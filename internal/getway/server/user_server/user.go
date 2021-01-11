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
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
)

func AccountLogin(ctx context.Context, login _map.AccountLogin) (*NUserPb.LoginResponse, *xerror.Err) {
	rep, err := xclient.NUserServer.AccountLogin(ctx, &NUserPb.UserLoginRequest{
		Account:  login.Account,
		Password: login.Password,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.EmptyData {
			e = e.SetMessage("账号或者密码错误")
		}
		return nil, e
	}
	return rep, nil
}
