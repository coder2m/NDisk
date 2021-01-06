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
	"github.com/myxy99/component/xlog"
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
			return nil, xrpc.NewError(xrpc.EmptyData)
		}
		xlog.Error("AccountLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return nil, err
	}
	if !u.CheckPassword(request.Password) {
		return nil, xrpc.NewError(xrpc.EmptyData)
	}
	token, err := accessToken.CreateAccessToken(ctx, xcast.ToUint64(u.ID))
	if !errors.Is(err, nil) {
		xlog.Error("AccountLogin", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return nil, err
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

func (s Server) SMSLogin(ctx context.Context, request *NUserPb.UserLoginRequest) (*NUserPb.LoginResponse, error) {
	panic("implement me")
}

func (s Server) UserRegister(ctx context.Context, request *NUserPb.UserRegisterRequest) (*NUserPb.NilResponse, error) {
	panic("implement me")
}

func (s Server) SendEmail(ctx context.Context, request *NUserPb.SendEmailRequest) (*NUserPb.NilResponse, error) {
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
