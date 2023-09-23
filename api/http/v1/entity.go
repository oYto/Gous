package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CodeSuccess           ErrCode = 0     // http 请求成功
	CodeBodyBindErr       ErrCode = 10001 // 参数解析错误
	CodeParamErr          ErrCode = 10002 // 请求参数不合法
	CodeRegisterErr       ErrCode = 10003 // 注册错误
	CodeLoginErr          ErrCode = 10003 // 登录错误
	CodeLogoutErr         ErrCode = 10004 // 登出错误
	CodeLogoffErr         ErrCode = 10004 // 注销错误
	CodeGetUserInfoErr    ErrCode = 10005 // 获取用户信息错误
	CodeUpdateUserInfoErr ErrCode = 10006 // 更新用户信息错误
)

type (
	DebugType int // debug 类型
	ErrCode   int // 错误码
)

// HttpResponse http 响应结构体，用户存储返回给客户端的数据
type HttpResponse struct {
	Code ErrCode     `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ResponseWithError 请求出错，根据错误码和信息返回错误给客户端
func (rsp *HttpResponse) ResponseWithError(c *gin.Context, code ErrCode, msg string) {
	rsp.Code = code
	rsp.Msg = msg
	c.JSON(http.StatusInternalServerError, rsp)
}

// ResponseSuccess 请求成功
func (rsp *HttpResponse) ResponseSuccess(c *gin.Context) {
	rsp.Code = CodeSuccess
	rsp.Msg = "success"
	c.JSON(http.StatusOK, rsp)
}

// ResponseWithData 带有数据的成功响应
func (rsp *HttpResponse) ResponseWithData(c *gin.Context, data interface{}) {
	rsp.Code = http.StatusOK
	rsp.Msg = "success"
	rsp.Data = data
	c.JSON(http.StatusOK, rsp)
}
