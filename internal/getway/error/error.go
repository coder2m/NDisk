/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 19:36
 **/
package xerror

import (
	"github.com/myxy99/component/pkg/xcode"
)

type Err struct {
	ErrorCode    uint32 `json:"err_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func (e *Err) Error() string {
	return e.ErrorMessage
}

func (e *Err) SetMessage(msg string) *Err {
	e.ErrorMessage = msg
	return e
}

func NewErr(code uint32, message string) *Err {
	return &Err{
		ErrorCode:    code,
		ErrorMessage: message,
	}
}

func NewErrRPC(err error) *Err {
	gst := xcode.ExtractCodes(err)
	if gst.Code <= 10000 {
		// todo
		gst=gst.SetMsg("server internal error")
	}
	return NewErr(gst.GetCodeAsUint32(), gst.Message)
}
