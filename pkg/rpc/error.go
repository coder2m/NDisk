/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 18:37
 **/
package xrpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	//业务错误 高于 10000
	EmptyData codes.Code = iota + 10000
)

var errDepot = map[codes.Code]string{
	EmptyData: "数据为找到",
}

func NewError(c codes.Code) error {
	return status.Error(c, errDepot[c])
}
