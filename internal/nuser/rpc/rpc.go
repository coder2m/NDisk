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
	"github.com/go-sql-driver/mysql"
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
	"github.com/myxy99/ndisk/internal/nuser/server/token"
	"github.com/myxy99/ndisk/pkg/constant"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrand "github.com/myxy99/ndisk/pkg/rand"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
	"gorm.io/gorm"
	"strings"
)

type Server struct{}

func (s Server) AccountLogin(ctx context.Context, request *NUserPb.UserLoginRequest) (rep *NUserPb.LoginResponse, err error) {
	var req = _map.AccountLogin{
		Account:  request.Account,
		Password: request.Password,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("accountLogin data validation error : %s", xvalidator.GetMsg(err).Error())
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
		return nil, xcode.BusinessCode(xrpc.AccountLoginErrCode)
	}
	if !u.CheckPassword(request.Password) {
		return nil, xcode.BusinessCode(xrpc.EmptyData)
	}
	if u.Status != 1 {
		return nil, xcode.BusinessCode(xrpc.LoginUserBanErrCode)
	}
	createAccessToken, err := xclient.RedisToken().CreateAccessToken(ctx, xcast.ToUint64(u.ID))
	if !errors.Is(err, nil) {
		xlog.Error("AccountLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return nil, xcode.BusinessCode(xrpc.AccountLoginErrCode)
	}
	return &NUserPb.LoginResponse{
		Info: &NUserPb.UserInfo{
			Uid:         xcast.ToUint64(u.ID),
			Name:        u.Name,
			Alias:       u.Alias,
			Email:       u.Email,
			Tel:         u.Tel,
			Status:      u.Status,
			EmailStatus: u.EmailStatus,
			CreatedAt:   xcast.ToUint64(u.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(u.UpdatedAt.Unix()),
		},
		Token: &NUserPb.Token{
			AccountToken: createAccessToken.AccessToken,
			RefreshToken: createAccessToken.RefreshToken,
		},
	}, nil
}

func (s Server) SMSSend(ctx context.Context, request *NUserPb.SendRequest) (nilR *NUserPb.NilResponse, err error) {
	var req = _map.Phone{
		Number: request.Account,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return nilR, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("SMSSend data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	if model.MainRedis().Exists(ctx, constant.SendVerificationCode.Format(request.Type, req.Number)).Val() > 0 {
		return nilR, xcode.BusinessCode(xrpc.FrequentOperationErrCode).SetMsgf("SMSSend frequent operation to phone:%v type:%+v", req.Number, request.Type)
	}
	code := xrand.CreateRandomNumber(constant.VerificationCodeLength)
	err = model.MainRedis().Set(ctx, constant.SendVerificationCode.Format(request.Type, req.Number), code, constant.VerificationEffectiveTime).Err()
	if !errors.Is(err, nil) {
		xlog.Error("SMSSend", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("redis"))
		return nilR, xcode.BusinessCode(xrpc.SMSSendErrCode)
	}
	smsRequest := xsms.SmsRequest{
		PhoneNumbers:  req.Number,
		TemplateParam: fmt.Sprintf(`{"code":"%s"}`, code),
	}
	res, err := xclient.SMSMain().Send(&smsRequest)
	if !errors.Is(err, nil) || !res.IsSuccess() {
		xlog.Error("SMSSend", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.Any("smsRequest", smsRequest))
		return nilR, xcode.BusinessCode(xrpc.SMSSendErrCode)
	}
	return new(NUserPb.NilResponse), nil
}

func (s Server) SMSLogin(ctx context.Context, request *NUserPb.SMSLoginRequest) (rep *NUserPb.LoginResponse, err error) {
	var req = _map.SMSLogin{
		Tel:  request.Tel,
		Code: request.Code,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("SMSSend data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	code := model.MainRedis().Get(ctx, constant.SendVerificationCode.Format(NUserPb.ActionType_Login_Type, req.Tel)).String()
	if code != req.Code {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("code Mismatch")
	}
	model.MainRedis().Del(ctx, constant.SendVerificationCode.Format(NUserPb.ActionType_Login_Type, req.Tel))
	u := new(user.User)
	err = u.GetByWhere(ctx, map[string][]interface{}{
		"tel =?": {req.Tel},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		xlog.Error("SMSLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.SMSLoginErrCode)
	}
	if u.Status != 1 {
		return nil, xcode.BusinessCode(xrpc.LoginUserBanErrCode)
	}
	createAccessToken, err := xclient.RedisToken().CreateAccessToken(ctx, xcast.ToUint64(u.ID))
	if !errors.Is(err, nil) {
		xlog.Error("SMSLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return nil, xcode.BusinessCode(xrpc.SMSLoginErrCode)
	}
	return &NUserPb.LoginResponse{
		Info: &NUserPb.UserInfo{
			Uid:         xcast.ToUint64(u.ID),
			Name:        u.Name,
			Alias:       u.Alias,
			Email:       u.Email,
			Tel:         u.Tel,
			Status:      u.Status,
			EmailStatus: u.EmailStatus,
			CreatedAt:   xcast.ToUint64(u.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(u.UpdatedAt.Unix()),
		},
		Token: &NUserPb.Token{
			AccountToken: createAccessToken.AccessToken,
			RefreshToken: createAccessToken.RefreshToken,
		},
	}, nil
}

func (s Server) SendEmail(ctx context.Context, request *NUserPb.SendRequest) (rep *NUserPb.NilResponse, err error) {
	var req = _map.Email{
		Email: request.Account,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("SendEmail data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	if model.MainRedis().Exists(ctx, constant.SendVerificationCode.Format(request.Type, req.Email)).Val() > 0 {
		return rep, xcode.BusinessCode(xrpc.FrequentOperationErrCode).SetMsgf("SendEmail frequent operation to email:%v type:%+v", req.Email, request.Type)
	}
	code := xrand.CreateRandomString(constant.VerificationCodeLength)
	err = model.MainRedis().Set(ctx, constant.SendVerificationCode.Format(request.Type, req.Email), code, constant.VerificationEffectiveTime).Err()
	if !errors.Is(err, nil) {
		xlog.Error("SendEmail", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("redis"))
		return rep, xcode.BusinessCode(xrpc.SendEmailErrCode)
	}
	err = xclient.EmailMain().SendEmail([]string{req.Email}, "验证码", fmt.Sprintf("你的验证码是：%s", code))
	if !errors.Is(err, nil) {
		xlog.Error("SendEmail", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.SendEmailErrCode)
	}
	return new(NUserPb.NilResponse), nil
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
	model.MainRedis().Del(ctx, constant.SendVerificationCode.Format(NUserPb.ActionType_Register_Type, req.Tel))

	if ok := new(user.User).ExistWhere(ctx, map[string][]interface{}{
		"name = ? or tel =? or email=?": {req.Name, req.Tel, req.Email},
	}); ok {
		return rep, xcode.BusinessCode(xrpc.DataExistErrCode)
	}

	var u = &user.User{Name: req.Name, Alias: req.Alias, Tel: req.Tel, Email: req.Email, Password: req.Password}
	err = u.SetPassword()
	err = u.Add(ctx)
	if !errors.Is(err, nil) {
		xlog.Error("User Register", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return rep, xcode.BusinessCode(xrpc.UserRegisterErrCode)
	}
	return new(NUserPb.NilResponse), nil
}

func (s Server) RetrievePassword(ctx context.Context, request *NUserPb.RetrievePasswordRequest) (rep *NUserPb.NilResponse, err error) {
	var req = _map.RetrievePassword{
		Account:  request.Account,
		Password: request.Password,
		Code:     request.Code,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("RetrievePassword data validation error : %s", xvalidator.GetMsg(err).Error())
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
		return nil, xcode.BusinessCode(xrpc.RetrievePasswordErrCode)
	}
	u.Password = req.Password
	err = u.SetPassword()
	err = u.UpdateWhere(ctx, map[string][]interface{}{
		"id=?": {u.ID},
	}, "password", u.Password)
	if !errors.Is(err, nil) {
		xlog.Error("RetrievePassword", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.RetrievePasswordErrCode)
	}
	return new(NUserPb.NilResponse), nil
}

func (s Server) GetUserById(ctx context.Context, info *NUserPb.UserInfo) (rep *NUserPb.UserInfo, err error) {
	req := _map.Id{
		Id: xcast.ToUint(info.Uid),
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("GetUserById data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	u := &user.User{
		Model: new(gorm.Model),
	}
	u.ID = req.Id
	err = u.GetById(ctx, false)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rep, xcode.BusinessCode(xrpc.EmptyData)
		}
		xlog.Error("GetUserById", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return rep, xcode.BusinessCode(xrpc.GetUserByIdErrCode)
	}
	return &NUserPb.UserInfo{
		Uid:         xcast.ToUint64(u.ID),
		Name:        u.Name,
		Alias:       u.Alias,
		Tel:         u.Tel,
		Email:       u.Email,
		Status:      u.Status,
		EmailStatus: u.EmailStatus,
		CreatedAt:   xcast.ToUint64(u.CreatedAt.Unix()),
		UpdatedAt:   xcast.ToUint64(u.UpdatedAt.Unix()),
		DeletedAt:   xcast.ToUint64(u.DeletedAt.Time.Unix()),
	}, nil
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
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("GetUserList data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	var (
		data  []user.User
		where map[string][]interface{}
	)
	if req.Keyword != "" {
		where = map[string][]interface{}{
			"name like ? or alias like ? or tel like ? or email like ?": {
				"%" + req.Keyword + "%",
				"%" + req.Keyword + "%",
				"%" + req.Keyword + "%",
				"%" + req.Keyword + "%",
			},
		}
	}
	total, err := new(user.User).Get(ctx, xcast.ToInt(req.PageSize*(req.Page-1)), xcast.ToInt(req.PageSize), &data, where, req.IsDelete)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		xlog.Error("GetUserList", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.GetUserListErrCode)
	}
	var userList = make([]*NUserPb.UserInfo, len(data))
	for i, datum := range data {
		userList[i] = &NUserPb.UserInfo{
			Uid:         xcast.ToUint64(datum.ID),
			Name:        datum.Name,
			Alias:       datum.Alias,
			Tel:         datum.Tel,
			Email:       datum.Email,
			Status:      datum.Status,
			EmailStatus: datum.EmailStatus,
			CreatedAt:   xcast.ToUint64(datum.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(datum.UpdatedAt.Unix()),
			DeletedAt:   xcast.ToUint64(datum.DeletedAt.Time.Unix()),
		}
	}
	return &NUserPb.UserListResponse{
		List:  userList,
		Count: xcast.ToUint32(total),
	}, nil
}

func (s Server) UpdateUserStatus(ctx context.Context, info *NUserPb.UserInfo) (rep *NUserPb.NilResponse, err error) {
	req := _map.UpdateUserStatus{
		Uid:    info.Uid,
		Status: info.Status,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("UpdateUserStatus data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	u := &user.User{
		Model: new(gorm.Model),
	}
	u.ID = xcast.ToUint(req.Uid)
	err = u.UpdateStatus(ctx, req.Status)
	if !errors.Is(err, nil) && err != gorm.ErrRecordNotFound {
		xlog.Error("UpdateUserStatus", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.UpdateUserStatusErrCode)
	}
	return new(NUserPb.NilResponse), nil
}

func (s Server) UpdateUserEmailStatus(ctx context.Context, info *NUserPb.UserInfo) (rep *NUserPb.NilResponse, err error) {
	req := _map.UpdateUserStatus{
		Uid:    info.Uid,
		Status: info.EmailStatus,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("UpdateUserEmailStatus data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	u := &user.User{
		Model: new(gorm.Model),
	}
	u.ID = xcast.ToUint(req.Uid)
	err = u.UpdateEmailStatus(ctx, req.Status)
	if !errors.Is(err, nil) && err != gorm.ErrRecordNotFound {
		xlog.Error("UpdateUserEmailStatus", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.UpdateUserEmailStatusErrCode)
	}
	return new(NUserPb.NilResponse), nil
}

func (s Server) UpdateUser(ctx context.Context, info *NUserPb.UserInfo) (rep *NUserPb.NilResponse, err error) {
	req := _map.UpdateUser{
		Uid:      info.Uid,
		Name:     info.Name,
		Alias:    info.Alias,
		Email:    info.Email,
		Tel:      info.Tel,
		Password: info.Password,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("UpdateUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	u := &user.User{
		Name:     req.Name,
		Alias:    req.Alias,
		Tel:      req.Tel,
		Email:    req.Email,
		Password: req.Password,
	}
	if u.Password != "" {
		_ = u.SetPassword()
	}
	err = u.UpdatesWhere(ctx, map[string][]interface{}{"id=?": {req.Uid}})
	if !errors.Is(err, nil) && err != gorm.ErrRecordNotFound {
		xlog.Error("UpdateUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.UpdateUserErrCode)
	}
	return new(NUserPb.NilResponse), nil
}

func (s Server) DelUsers(ctx context.Context, list *NUserPb.UidList) (rep *NUserPb.ChangeNumResponse, err error) {
	if len(list.Uid) <= 0 {
		return rep, err
	}
	var data = make([]string, len(list.Uid))
	for i, u := range list.Uid {
		data[i] = xcast.ToString(u)
	}
	where := map[string][]interface{}{
		"id IN (?)": {strings.Join(data, ",")},
	}
	count, err := new(user.User).Del(ctx, where)
	if !errors.Is(err, nil) && err != gorm.ErrRecordNotFound {
		xlog.Error("DelUsers", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.DelUsersErrCode)
	}
	rep = new(NUserPb.ChangeNumResponse)
	rep.Count = xcast.ToUint32(count)
	return rep, nil
}

func (s Server) RecoverDelUsers(ctx context.Context, list *NUserPb.UidList) (rep *NUserPb.ChangeNumResponse, err error) {
	if len(list.Uid) <= 0 {
		return rep, err
	}
	var data = make([]string, len(list.Uid))
	for i, u := range list.Uid {
		data[i] = xcast.ToString(u)
	}
	where := map[string][]interface{}{
		"id IN (?)": {strings.Join(data, ",")},
	}
	count, err := new(user.User).DelRes(ctx, where)
	if !errors.Is(err, nil) && err != gorm.ErrRecordNotFound {
		xlog.Error("RecoverDelUsers", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.RecoverDelUsersErrCode)
	}
	rep = new(NUserPb.ChangeNumResponse)
	rep.Count = xcast.ToUint32(count)
	return rep, nil
}

func (s Server) CreateUsers(ctx context.Context, list *NUserPb.UserList) (rep *NUserPb.ChangeNumResponse, err error) {
	if len(list.List) <= 0 {
		return rep, nil
	}
	if len(list.List) > 200 {
		return rep, xcode.BusinessCode(xrpc.MaximumNumberErrCode)
	}
	var data = make([]user.User, len(list.List))
	for i, info := range list.List {
		data[i] = user.User{
			Name:     info.Name,
			Alias:    info.Alias,
			Tel:      info.Tel,
			Email:    info.Email,
			Password: info.Password,
		}
		_ = data[i].SetPassword()
	}
	count, err := new(user.User).Adds(ctx, data)
	if !errors.Is(err, nil) && err != gorm.ErrRecordNotFound {
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Number == 1062 { // Duplicate
				return nil, xcode.BusinessCode(xrpc.CreateUsersErrCode).SetMsg("数据已经存在")
			}
		}
		xlog.Error("CreateUsers", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.CreateUsersErrCode)
	}
	rep = new(NUserPb.ChangeNumResponse)
	rep.Count = xcast.ToUint32(count)
	return rep, nil
}

func (s Server) VerifyUsers(ctx context.Context, r *NUserPb.Token) (rep *NUserPb.UserInfo, err error) {
	req := _map.UserToken{
		Token: r.AccountToken,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("VerifyUsers data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	i, _ := xclient.RedisToken().DecoderAccessToken(ctx, req.Token)
	if i.Uid <= 0 || i.Type != token.AccessTokenType {
		return rep, xcode.BusinessCode(xrpc.VerifyUsersTokenErrCode)
	}
	return s.GetUserById(ctx, &NUserPb.UserInfo{Uid: i.Uid})
}

func (s Server) RefreshToken(ctx context.Context, r *NUserPb.Token) (rep *NUserPb.Token, err error) {
	req := _map.UserToken{
		Token: r.RefreshToken,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("RefreshToken data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	refreshAccessToken, err := xclient.RedisToken().RefreshAccessToken(ctx, req.Token)
	if !errors.Is(err, nil) {
		xlog.Error("RefreshToken", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return nil, xcode.BusinessCode(xrpc.RefreshTokenErrCode)
	}
	return &NUserPb.Token{
		AccountToken: refreshAccessToken.AccessToken,
		RefreshToken: refreshAccessToken.RefreshToken,
	}, nil
}

func (s Server) CheckCode(ctx context.Context, r *NUserPb.CheckCodeRequest) (rep *NUserPb.NilResponse, err error) {
	var req = _map.CheckCode{
		Account: r.Account,
		Code:    r.Code,
		Type:    r.Type,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("CheckCode data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	code := model.MainRedis().Get(ctx, constant.SendVerificationCode.Format(req.Type, req.Account)).String()
	if code != req.Code {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("code Mismatch")
	}
	model.MainRedis().Del(ctx, constant.SendVerificationCode.Format(req.Type, req.Account))
	return new(NUserPb.NilResponse), nil
}

func (s Server) GetUserListByUid(ctx context.Context, req *NUserPb.UidList) (rep *NUserPb.UserListResponse, err error) {
	if len(req.Uid) <= 0 {
		return rep, err
	}
	var data = make([]string, len(req.Uid))
	for i, u := range req.Uid {
		data[i] = xcast.ToString(u)
	}
	where := map[string][]interface{}{
		"id IN (?)": {strings.Join(data, ",")},
	}
	var userList []user.User
	err = new(user.User).GetAll(ctx, &userList, where)
	if !errors.Is(err, nil) && err != gorm.ErrRecordNotFound {
		xlog.Error("GetUserListByUid", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, xcode.BusinessCode(xrpc.GetUserListByUidErrCode)
	}
	var list = make([]*NUserPb.UserInfo, len(userList))
	for i, datum := range userList {
		list[i] = &NUserPb.UserInfo{
			Uid:         xcast.ToUint64(datum.ID),
			Name:        datum.Name,
			Alias:       datum.Alias,
			Tel:         datum.Tel,
			Email:       datum.Email,
			Status:      datum.Status,
			EmailStatus: datum.EmailStatus,
			CreatedAt:   xcast.ToUint64(datum.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(datum.UpdatedAt.Unix()),
			DeletedAt:   xcast.ToUint64(datum.DeletedAt.Time.Unix()),
		}
	}
	return &NUserPb.UserListResponse{
		List:  list,
		Count: xcast.ToUint32(len(userList)),
	}, nil
}
