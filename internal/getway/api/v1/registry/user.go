/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 16:23
 **/
package registry

import (
	nh "github.com/myxy99/ndisk/internal/getway/api/v1/handle/user"
	"github.com/myxy99/ndisk/internal/getway/api/v1/middleware"
)

func init() {
	user := V1().Group("/auth")
	user.POST("/login", middleware.Recaptcha(), nh.AccountLogin)
}
