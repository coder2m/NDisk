/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 19:21
 **/
package auth_server

import (
	"context"
	"errors"
	xclient "github.com/myxy99/ndisk/internal/getway/client"
	xerror "github.com/myxy99/ndisk/internal/getway/error"
	_map "github.com/myxy99/ndisk/internal/getway/map"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
)

func AccountLogin(ctx context.Context, login _map.AccountLogin) (map[string]interface{}, *xerror.Err) {
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
	return map[string]interface{}{
		"info": _map.UserInfo{
			Uid:         rep.Info.Uid,
			Name:        rep.Info.Name,
			Alias:       rep.Info.Alias,
			Tel:         rep.Info.Tel,
			Email:       rep.Info.Email,
			Status:      rep.Info.Status,
			EmailStatus: rep.Info.EmailStatus,
			CreatedAt:   rep.Info.CreatedAt,
			UpdatedAt:   rep.Info.UpdatedAt,
		},
		"token": _map.Token{
			AccountToken: rep.Token.AccountToken,
			RefreshToken: rep.Token.RefreshToken,
		},
	}, nil
}

func SMSSend(ctx context.Context, to _map.SMSSend) *xerror.Err {
	_, err := xclient.NUserServer.SMSSend(ctx, &NUserPb.SendRequest{
		Account: to.Tel,
		Type:    to.Type,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.FrequentOperationErrCode {
			e = e.SetMessage("操作频繁")
		}
		return e
	}
	return nil
}

func SMSLogin(ctx context.Context, to _map.SMSLogin) (map[string]interface{}, *xerror.Err) {
	rep, err := xclient.NUserServer.SMSLogin(ctx, &NUserPb.SMSLoginRequest{
		Tel:  to.Tel,
		Code: to.Code,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.ValidationErrCode {
			e = e.SetMessage("电话号码或验证码错误")
		}
		return nil, e
	}
	return map[string]interface{}{
		"info": _map.UserInfo{
			Uid:         rep.Info.Uid,
			Name:        rep.Info.Name,
			Alias:       rep.Info.Alias,
			Tel:         rep.Info.Tel,
			Email:       rep.Info.Email,
			EmailStatus: rep.Info.EmailStatus,
			Status:      rep.Info.Status,
			CreatedAt:   rep.Info.CreatedAt,
			UpdatedAt:   rep.Info.UpdatedAt,
		},
		"token": _map.Token{
			AccountToken: rep.Token.AccountToken,
			RefreshToken: rep.Token.RefreshToken,
		},
	}, nil
}

func Register(ctx context.Context, r _map.UserRegister) *xerror.Err {
	_, err := xclient.NUserServer.UserRegister(ctx, &NUserPb.UserRegisterRequest{
		Info: &NUserPb.UserInfo{
			Name:     r.Name,
			Alias:    r.Alias,
			Tel:      r.Tel,
			Email:    r.Email,
			Password: r.Password,
		},
		Code: r.Code,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.ValidationErrCode {
			e = e.SetMessage("参数或验证码错误")
		} else if e.ErrorCode == xrpc.DataExistErrCode {
			e = e.SetMessage("用户名或邮箱或电话占用")
		}
		return e
	}
	return nil
}

func SendEmail(ctx context.Context, r _map.EmailSend) *xerror.Err {
	_, err := xclient.NUserServer.SendEmail(ctx, &NUserPb.SendRequest{
		Account: r.Email,
		Type:    r.Type,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.FrequentOperationErrCode {
			e = e.SetMessage("操作频繁")
		}
		return e
	}
	return nil
}

func Retrieve(ctx context.Context, r _map.RetrievePassword) *xerror.Err {
	_, err := xclient.NUserServer.RetrievePassword(ctx, &NUserPb.RetrievePasswordRequest{
		Account:  r.Account,
		Password: r.Password,
		Code:     r.Code,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.ValidationErrCode || e.ErrorCode == xrpc.EmptyData {
			e = e.SetMessage("账号或者验证码错误")
		}
		return e
	}
	return nil
}

func RefreshToken(ctx context.Context, r _map.UserToken) (interface{}, *xerror.Err) {
	data, err := xclient.NUserServer.RefreshToken(ctx, &NUserPb.Token{
		AccountToken: r.Token,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		return nil, e
	}
	return data, nil
}

func BindEmail(ctx context.Context, r _map.BindEmail) *xerror.Err {
	_, err := xclient.NUserServer.CheckCode(ctx, &NUserPb.CheckCodeRequest{
		Account: r.Email,
		Code:    r.Code,
		Type:    NUserPb.ActionType_EmailAttest_Type,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.ValidationErrCode {
			e = e.SetMessage("账号或者验证码错误")
		}
		return e
	}
	_, err = xclient.NUserServer.UpdateUserEmailStatus(ctx, &NUserPb.UserInfo{
		Uid:         r.Uid,
		EmailStatus: 1,
	})

	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		return e
	}
	return nil
}
