/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 17:33
 **/
package rpc

import (
	"context"
	"errors"
	"github.com/myxy99/ndisk/internal/nuser/model/user"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
	"gorm.io/gorm"
)

type Server struct{}

func (s Server) AccountLogin(ctx context.Context, request *NUserPb.UserLoginRequest) (*NUserPb.LoginResponse, error) {
	u := new(user.User)
	err := u.GetByWhere(map[string][]interface{}{
		"name = ? or tel =? or email=?": {request.Account, request.Account, request.Account},
	})
	if errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xrpc.NewError(xrpc.EmptyData)
		}
		return nil, err
	}
	if !u.CheckPassword(request.Password) {
		return nil, xrpc.NewError(xrpc.EmptyData)
	}

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
