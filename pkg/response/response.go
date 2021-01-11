/**
 * @Author: myxy99 <myxy99@foxmail.com>
 * @Date: 2020-10-14 14:28
 */
package R

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

type Res struct {
	StatusCode int         `json:"status_code,omitempty"`
	Msg        string      `json:"msg,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

//R.Ok(c, "自定义msg",data)
func Ok(c *gin.Context, data interface{}) {
	Response(c, Success, MSG_OK, data, http.StatusOK)
}

//R.Error(c, "自定义msg",data)
func Error(c *gin.Context, data interface{}) {
	Response(c, File, MSG_ERR, data, http.StatusOK)
}

//R.Response(c,1,"msg",data,200)
func Response(c *gin.Context, code int, msg string, data interface{}, status int) {
	c.Render(status, render.JSON{Data: Res{
		code,
		msg,
		data,
	}})
}
