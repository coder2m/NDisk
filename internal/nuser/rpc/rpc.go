/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 17:33
 **/
package rpc

import (
	"context"
	"errors"
	"fmt"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xcast"
	"github.com/myxy99/component/pkg/xcode"
	"github.com/myxy99/component/pkg/xvalidator"
	xsms "github.com/myxy99/component/xinvoker/sms"
	"github.com/myxy99/component/xlog"
	xclient "github.com/myxy99/ndisk/internal/nuser/client"
	_map "github.com/myxy99/ndisk/internal/nuser/map"
	"github.com/myxy99/ndisk/internal/nuser/model"
	"github.com/myxy99/ndisk/internal/nuser/model/user"
	redisToken "github.com/myxy99/ndisk/internal/nuser/server/token/redis"
	"github.com/myxy99/ndisk/pkg/constant"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrand "github.com/myxy99/ndisk/pkg/rand"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
	"gorm.io/gorm"
)

var (
	accessToken = redisToken.NewAccessToken(model.MainRedis())
)

type Server struct{}

func (s Server) AccountLogin(ctx context.Context, request *NUserPb.UserLoginRequest) (rep *NUserPb.LoginResponse, err error) {
	var req = _map.AccountLogin{
		Account:  request.Account,
		Password: request.Password,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("accountLogin data validation error : %s", msg.Error())
	}
	u := new(user.User)
	err = u.GetByWhere(ctx, map[string][]interface{}{
		"name = ? or tel =? or email=?": {request.Account, request.Account, request.Account},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		xlog.Error("AccountLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
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
			Tel:       u.Tel,
			Status:    u.Status,
			CreatedAt: xcast.ToUint64(u.CreatedAt.Unix()),
			UpdatedAt: xcast.ToUint64(u.UpdatedAt.Unix()),
		},
		Token: &NUserPb.Token{
			AccountToken: token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, xcode.OK
}

func (s Server) SMSSend(ctx context.Context, request *NUserPb.SendRequest) (nilR *NUserPb.NilResponse, err error) {
	var req = _map.Phone{
		Number: request.Account,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return nilR, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("SMSSend data validation error : %s", msg.Error())
	}
	if model.MainRedis().Exists(ctx, constant.SendVerificationCode.Format(request.Type, req.Number)).Val() > 0 {
		return nilR, xcode.BusinessCode(xrpc.FrequentOperationErrCode).SetMsgf("SMSSend frequent operation to phone:%v type:", req.Number, request.Type)
	}
	code := xrand.CreateRandomNumber(constant.VerificationCodeLength)
	err = model.MainRedis().Set(ctx, constant.SendVerificationCode.Format(request.Type, req.Number), code, constant.VerificationEffectiveTime).Err()
	if !errors.Is(err, nil) {
		xlog.Error("SMSSend", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("redis"))
		return nilR, xcode.BusinessCode(xrpc.SMSSendErrCode).SetMsgf("SMSSend Send error : %s", err.Error())
	}
	smsRequest := xsms.SmsRequest{
		PhoneNumbers:  req.Number,
		TemplateParam: fmt.Sprintf(`{"code":"%s"}`, code),
	}
	res, err := xclient.SMSMain().Send(&smsRequest)
	if !errors.Is(err, nil) || !res.IsSuccess() {
		xlog.Error("SMSSend", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.Any("smsRequest", smsRequest))
		return nilR, xcode.BusinessCode(xrpc.SMSSendErrCode).SetMsgf("SMSSend Send error : %s", err.Error())
	}
	return nilR, xcode.OK
}

func (s Server) SMSLogin(ctx context.Context, request *NUserPb.SMSLoginRequest) (rep *NUserPb.LoginResponse, err error) {
	var req = _map.SMSLogin{
		Tel:  request.Tel,
		Code: request.Code,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("SMSSend data validation error : %s", msg.Error())
	}
	code := model.MainRedis().Get(ctx, constant.SendVerificationCode.Format(NUserPb.ActionType_Login_Type, req.Tel)).String()
	if code != req.Code {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("code Mismatch")
	}
	u := new(user.User)
	err = u.GetByWhere(ctx, map[string][]interface{}{
		"tel =?": {req.Tel},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		xlog.Error("SMSLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.SMSLoginErrCode).SetMsgf("SMSLogin error : %s", err)
	}
	token, err := accessToken.CreateAccessToken(ctx, xcast.ToUint64(u.ID))
	if !errors.Is(err, nil) {
		xlog.Error("SMSLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return nil, xcode.BusinessCode(xrpc.SMSLoginErrCode).SetMsgf("SMSLogin error : %s", err)
	}
	return &NUserPb.LoginResponse{
		Info: &NUserPb.UserInfo{
			Uid:       xcast.ToUint64(u.ID),
			Name:      u.Name,
			Alias:     u.Alias,
			Email:     u.Email,
			Tel:       u.Tel,
			Status:    u.Status,
			CreatedAt: xcast.ToUint64(u.CreatedAt.Unix()),
			UpdatedAt: xcast.ToUint64(u.UpdatedAt.Unix()),
		},
		Token: &NUserPb.Token{
			AccountToken: token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, xcode.OK
}

func (s Server) SendEmail(ctx context.Context, request *NUserPb.SendRequest) (rep *NUserPb.NilResponse, err error) {
	var req = _map.Email{
		Email: request.Account,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("SendEmail data validation error : %s", msg.Error())
	}
	if model.MainRedis().Exists(ctx, constant.SendVerificationCode.Format(request.Type, req.Email)).Val() > 0 {
		return rep, xcode.BusinessCode(xrpc.FrequentOperationErrCode).SetMsgf("SendEmail frequent operation to email:%v type:", req.Email, request.Type)
	}
	code := xrand.CreateRandomString(constant.VerificationCodeLength)
	err = model.MainRedis().Set(ctx, constant.SendVerificationCode.Format(request.Type, req.Email), code, constant.VerificationEffectiveTime).Err()
	if !errors.Is(err, nil) {
		xlog.Error("SendEmail", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("redis"))
		return rep, xcode.BusinessCode(xrpc.SendEmailErrCode).SetMsgf("SendEmail Send error : %s", err.Error())
	}
	err = xclient.EmailMain().SendEmail([]string{req.Email}, "验证码", fmt.Sprintf("你的验证码是：%v", code))
	if !errors.Is(err, nil) {
		xlog.Error("SendEmail", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.SendEmailErrCode).SetMsgf("SendEmail  error : %s", err.Error())
	}
	return rep, xcode.OK
}

func (s Server) UserRegister(ctx context.Context, request *NUserPb.UserRegisterRequest) (rep *NUserPb.NilResponse, err error) {
	var req = _map.UserRegister{
		Name:     request.Info.Name,
		Alias:    request.Info.Alias,
		Email:    request.Info.Email,
		Tel:      request.Info.Tel,
		Password: request.Info.Password,
		Code:     request.Code}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("UserRegister data validation error : %s", msg.Error())
	}
	code := model.MainRedis().Get(ctx, constant.SendVerificationCode.Format(NUserPb.ActionType_Register_Type, req.Tel)).String()
	if code != req.Code {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("code Mismatch")
	}
	var u = &user.User{Name: req.Name, Alias: req.Alias, Tel: req.Tel, Email: req.Email, Password: req.Password}
	err = u.SetPassword()
	err = u.Add(ctx)
	if !errors.Is(err, nil) {
		xlog.Error("User Register", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return rep, xcode.BusinessCode(xrpc.UserRegisterErrCode).SetMsgf("UserRegister error : %s", err.Error())
	}
	return rep, xcode.OK
}

func (s Server) RetrievePassword(ctx context.Context, request *NUserPb.RetrievePasswordRequest) (rep *NUserPb.NilResponse, err error) {
	var req = _map.RetrievePassword{
		Account:  request.Account,
		Password: request.Password,
		Code:     request.Code,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("RetrievePassword data validation error : %s", msg.Error())
	}
	code := model.MainRedis().Get(ctx, constant.SendVerificationCode.Format(NUserPb.ActionType_Retrieve_Type, req.Account)).String()
	if code != req.Code {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("code Mismatch")
	}
	model.MainRedis().Del(ctx, constant.SendVerificationCode.Format(NUserPb.ActionType_Retrieve_Type, req.Account))
	u := new(user.User)
	err = u.GetByWhere(ctx, map[string][]interface{}{
		"tel =? or email=?": {req.Account, req.Account},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		xlog.Error("RetrievePassword", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.RetrievePasswordErrCode).SetMsgf("RetrievePassword error : %s", err)
	}
	u.Password = req.Password
	err = u.SetPassword()
	err = u.UpdateWhere(ctx, map[string][]interface{}{
		"id=?": {u.ID},
	}, "password", u.Password)
	if !errors.Is(err, nil) {
		xlog.Error("RetrievePassword", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.RetrievePasswordErrCode).SetMsgf("RetrievePassword error : %s", err)
	}
	return rep, xcode.OK
}

func (s Server) GetUserById(ctx context.Context, info *NUserPb.UserInfo) (rep *NUserPb.UserInfo, err error) {
	req := _map.IdMap{
		Id: xcast.ToUint(info.Uid),
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("GetUserById data validation error : %s", msg.Error())
	}
	u := new(user.User)
	u.ID = req.Id
	err = u.GetById(ctx, false)
	if !errors.Is(err, nil) {
		xlog.Error("GetUserById", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.GetUserByIdErrCode).SetMsgf("GetUserById error : %s", err)
	}
	return &NUserPb.UserInfo{
		Uid:       xcast.ToUint64(u.ID),
		Name:      u.Name,
		Alias:     u.Alias,
		Tel:       u.Tel,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: xcast.ToUint64(u.CreatedAt.Unix()),
		UpdatedAt: xcast.ToUint64(u.UpdatedAt.Unix()),
		DeletedAt: xcast.ToUint64(u.DeletedAt.Time.Unix()),
	}, err
}

func (s Server) GetUserList(ctx context.Context, request *NUserPb.PageRequest) (rep *NUserPb.UserListResponse, err error) {
	req := _map.PageList{
		Page:     request.Page,
		PageSize: request.Limit,
		Keyword:  request.Keyword,
		IsDelete: request.IsDelete,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		msg := xvalidator.GetMsg(err)
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("GetUserList data validation error : %s", msg.Error())
	}
	var data []user.User
	where := map[string][]interface{}{
		"name=? or alias=? or tel=? or email=?": {req.Keyword, req.Keyword, req.Keyword, req.Keyword},
	}
	total, err := new(user.User).Get(ctx, xcast.ToInt(req.PageSize*(req.Page-1)), xcast.ToInt(req.PageSize), &data, where, req.IsDelete)
	if !errors.Is(err, nil) {
		xlog.Error("GetUserList", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.GetUserListErrCode).SetMsgf("GetUserList error : %s", err)
	}
	var userList = make([]*NUserPb.UserInfo, len(data))
	for i, datum := range data {
		userList[i] = &NUserPb.UserInfo{
			Uid:       xcast.ToUint64(datum.ID),
			Name:      datum.Name,
			Alias:     datum.Alias,
			Tel:       datum.Tel,
			Email:     datum.Email,
			Status:    datum.Status,
			CreatedAt: xcast.ToUint64(datum.CreatedAt.Unix()),
			UpdatedAt: xcast.ToUint64(datum.UpdatedAt.Unix()),
			DeletedAt: xcast.ToUint64(datum.DeletedAt.Time.Unix()),
		}
	}
	return &NUserPb.UserListResponse{
		List:  userList,
		Count: xcast.ToUint32(total),
	}, err
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
