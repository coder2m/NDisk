/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 19:21
 **/
package auth_server

import (
	"context"
	"errors"

	"github.com/myxy99/component/pkg/xcast"
	xclient "github.com/myxy99/ndisk/internal/getway/client"
	xerror "github.com/myxy99/ndisk/internal/getway/error"
	_map "github.com/myxy99/ndisk/internal/getway/map"
	AuthorityPb "github.com/myxy99/ndisk/pkg/pb/authority"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
)

func AccountLogin(ctx context.Context, login _map.AccountLogin) (interface{}, *xerror.Err) {
	rep, err := xclient.NUserServer.AccountLogin(ctx, &NUserPb.UserLoginRequest{
		Account:  login.Account,
		Password: login.Password,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.EmptyData {
			e = e.SetMessage("账号或者密码错误")
		} else if e.ErrorCode == xrpc.LoginUserBanErrCode {
			e = e.SetMessage("用户被禁用")
		}
		return nil, e
	}
	return _map.Token{
		AccountToken: rep.Token.AccountToken,
		RefreshToken: rep.Token.RefreshToken,
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

func SMSLogin(ctx context.Context, to _map.SMSLogin) (interface{}, *xerror.Err) {
	rep, err := xclient.NUserServer.SMSLogin(ctx, &NUserPb.SMSLoginRequest{
		Tel:  to.Tel,
		Code: to.Code,
	})
	if !errors.Is(err, nil) {
		e := xerror.NewErrRPC(err)
		if e.ErrorCode == xrpc.ValidationErrCode {
			e = e.SetMessage("电话号码或验证码错误")
		} else if e.ErrorCode == xrpc.LoginUserBanErrCode {
			e = e.SetMessage("用户被禁用")
		}
		return nil, e
	}
	return _map.Token{
		AccountToken: rep.Token.AccountToken,
		RefreshToken: rep.Token.RefreshToken,
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

func GetPermissionAndMenuByRoles(ctx context.Context, roles []string) (map[string]interface{}, *xerror.Err) {
	data := make(map[string]interface{})
	for _, role := range roles {
		res, _ := xclient.AuthorityServer.GetPermissionAndMenuByRoles(ctx, &AuthorityPb.Target{
			To: role,
		})
		var menusList = make([]_map.MenuInfoRes, len(res.Menus))
		for i, menu := range res.Menus {
			menusList[i] = _map.MenuInfoRes{
				Id:          xcast.ToUint32(menu.Id),
				ParentId:    xcast.ToUint32(menu.ParentId),
				Path:        menu.Path,
				Name:        menu.Name,
				Description: menu.Description,
				IconClass:   menu.IconClass,
			}
		}
		data[role] = _map.RolesInfoRes{
			Name:  res.Name,
			Menus: menusList,
		}
	}
	return data, nil
}
