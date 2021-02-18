package R

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleInternalError(ctx *gin.Context) {
	Response(ctx, InternalError, "Internal Error", nil, http.StatusInternalServerError)
	return
}

func HandleBadRequest(ctx *gin.Context, data interface{}) {
	Response(ctx, BadRequest, "Bad Request", data, http.StatusBadRequest)
	return
}

func HandleNotFound(ctx *gin.Context) {
	Response(ctx, NotFound, "Not Found", nil, http.StatusNotFound)
	return
}

func HandleForbidden(ctx *gin.Context) {
	Response(ctx, Forbidden, "Forbidden", nil, http.StatusForbidden)
	return
}

func HandleCaptchaError(ctx *gin.Context) {
	Response(ctx, BadRequest, "captcha error", nil, http.StatusBadRequest)
	return
}
