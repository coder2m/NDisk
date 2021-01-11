/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 19:36
 **/
package xerror

import "github.com/myxy99/component/pkg/xcode"

type Err struct {
	Code    uint32
	Message string
}

func NewErr(code uint32, message string) *Err {
	return &Err{
		Code:    code,
		Message: message,
	}
}

func NewErrRPC(err error) *Err {
	gst := xcode.ExtractCodes(err)
	return NewErr(gst.GetCodeAsUint32(), gst.Message)
}
