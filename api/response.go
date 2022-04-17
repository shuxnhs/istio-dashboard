package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	ErrCode int         `json:"err_code"`
	Data    interface{} `json:"data"`
}

func NewResponse(code, errcode int, message string, data interface{}) *Result {
	return &Result{
		Code:    code,
		Message: message,
		ErrCode: errcode,
		Data:    data,
	}
}

func (r *Result) Resp(ctx *gin.Context) {
	ctx.JSON(r.Code, gin.H{
		"err_code": r.ErrCode,
		"message":  r.Message,
		"data":     r.Data,
	})
}

func ResponseOK(ctx *gin.Context) {
	NewResponse(http.StatusOK, CodeSuccess, "", nil).Resp(ctx)
}

func ResponseData(ctx *gin.Context, errCode int, data interface{}) {
	NewResponse(http.StatusOK, errCode, "", data).Resp(ctx)
}

func Response(ctx *gin.Context, code, errCode int, message string, data interface{}) {
	NewResponse(code, errCode, message, data).Resp(ctx)
}
