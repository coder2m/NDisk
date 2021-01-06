/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 17:33
 **/
package rpc

import (
	"context"
	"errors"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xcast"
	"github.com/myxy99/component/pkg/xcode"
	"github.com/myxy99/component/pkg/xvalidator"
	xsms "github.com/myxy99/component/xinvoker/sms"
	"github.com/myxy99/component/xlog"
	"github.com/myxy99/ndisk/internal/nuser/constant"
	_map "github.com/myxy99/ndisk/internal/nuser/map"
	"github.com/myxy99/ndisk/internal/nuser/model"
	"github.com/myxy99/ndisk/internal/nuser/model/user"
	redisToken "github.com/myxy99/ndisk/internal/nuser/server/token/redis"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
	"gorm.io/gorm"
)

var (
	accessToken = redisToken.NewAccessToken(model.MainRedis())
)

type Server struct{}

func (s Server) AccountLogin(ctx context.Context, request *NUserPb.UserLoginRequest) (*NUserPb.LoginResponse, error) {
	u := new(user.User)
	err := u.GetByWhere(map[string][]interface{}{
		"name = ? or tel =? or email=?": {request.Account, request.Account, request.Account},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		xlog.Error("AccountLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return nil, xcode.BusinessCode(xrpc.AccountLoginErrCode).SetMsgf("AccountLogin error : %s", err)
	}
	if !u.CheckPassword(request.Password) {
		return nil, xcode.BusinessCode(xrpc.EmptyData)
	}
	token, err := accessToken.CreateAccessToken(ctx, xcast.ToUint64(u.ID))
	if !errors.Is(err, nil) {
		xlog.Error("AccountLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return nil, xcode.BusinessCode(xrpc.AccountLoginErrCode).SetMsgf("AccountLogin error : %s", err)
	}
	return &NUserPb.LoginResponse{
		Info: &NUserPb.UserInfo{
			Uid:       xcast.ToUint64(u.ID),
			Name:      u.Name,
			Alias:     u.Alias,
			Email:     u.Email,
			Status:    u.Status,
			CreatedAt: xcast.ToUint64(u.CreatedAt.Unix()),
			UpdatedAt: xcast.ToUint64(u.UpdatedAt.Unix()),
		},
		Token: &NUserPb.Token{
			AccountToken: token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, err
}

func (s Server) SMSSend(ctx context.Context, request *NUserPb.SendRequest) (nilR *NUserPb.NilResponse, err error) {
	var phone _map.Phone
	phone.Number = request.Account
	err = xvalidator.Struct(phone)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return nilR, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("SMSSend data validation error : %s", msg.Error())
	}
	//todo 验证验证码存在
	if model.MainRedis().Exists(ctx, constant.SendSMS.Format(request.Type, request.Account)).Val() > 0 {
		return nilR, xcode.BusinessCode(xrpc.FrequentOperationErrCode).SetMsgf("SMSSend frequent operation to phone:%v type:", phone.Number, request.Type)
	}
	smsRequest := xsms.SmsRequest{
		PhoneNumbers:  phone.Number,
		TemplateParam: "{\"code\":\"2312\"}",
	}
	res, err := xsms.Invoker("main").Send(&smsRequest)
	if !errors.Is(err, nil) || !res.IsSuccess() {
		xlog.Error("SMSSend", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.Any("smsRequest", smsRequest))
	}
	return new(NUserPb.NilResponse), err
}

func (s Server) SMSLogin(ctx context.Context, request *NUserPb.SMSLoginRequest) (*NUserPb.LoginResponse, error) {
	panic("implement me")
}

func (s Server) SendEmail(ctx context.Context, request *NUserPb.SendRequest) (*NUserPb.NilResponse, error) {
	panic("implement me")
}

func (s Server) UserRegister(ctx context.Context, request *NUserPb.UserRegisterRequest) (*NUserPb.NilResponse, error) {
	panic("implement me")
}

func (s Server) RetrievePassword(ctx context.Context, request *NUserPb.RetrievePasswordRequest) (*NUserPb.NilResponse, error) {
	panic("implement me")
}

func (s Server) GetUserById(ctx context.Context, info *NUserPb.UserInfo) (*NUserPb.UserInfo, error) {
	panic("implement me")
}

func (s Server) GetUserList(ctx context.Context, request *NUserPb.PageRequest) (*NUserPb.UserListResponse, error) {
	panic("implement me")
}

func (s Server) UpdateUserStatus(ctx context.Context, info *NUserPb.UserInfo) (*NUserPb.NilResponse, error) {
	panic("implement me")
}

func (s Server) UpdateUser(ctx context.Context, info *NUserPb.UserInfo) (*NUserPb.NilResponse, error) {
	panic("implement me")
}

func (s Server) DelUsers(ctx context.Context, list *NUserPb.UidList) (*NUserPb.BatchOperationResponse, error) {
	panic("implement me")
}

func (s Server) RecoverDelUsers(ctx context.Context, list *NUserPb.UidList) (*NUserPb.BatchOperationResponse, error) {
	panic("implement me")
}

func (s Server) CreateUsers(ctx context.Context, list *NUserPb.UserList) (*NUserPb.BatchOperationResponse, error) {
	panic("implement me")
}

func (s Server) VerifyUsers(ctx context.Context, list *NUserPb.Token) (*NUserPb.UserInfo, error) {
	panic("implement me")
}
