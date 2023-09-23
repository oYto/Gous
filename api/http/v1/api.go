package v1

import (
	"Gous/config"
	"Gous/internal/service"
	"Gous/internal/utils"
	"Gous/pkg/constant"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Ping 健康检查
func Ping(c *gin.Context) {
	appConfig := config.GetGlobalConf().AppConfig
	conInfo, _ := json.MarshalIndent(appConfig, "", "")
	appInfo := fmt.Sprintf("Gous：\n\napp_name: %s\nversion: %s\n\n%s", appConfig.AppName, appConfig.Version, string(conInfo))
	c.String(http.StatusOK, appInfo)
}

// Register 注册
func Register(c *gin.Context) {
	// 请求
	req := &service.RegisterRequest{}
	// 响应
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(&req)
	if err != nil { // 参数解析错误
		log.Error("Gous：bind request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
	}

	// 开始注册
	if err := service.Register(req); err != nil {
		rsp.ResponseWithError(c, CodeRegisterErr, err.Error())
		return
	}
	rsp.ResponseSuccess(c)
}

// Login 登录
func Login(c *gin.Context) {
	req := &service.LoginRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request json err:%v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
	}

	// 获取 uuid
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx := context.WithValue(context.Background(), "uuid", uuid)
	log.Infof("loggin start, user:%s, password:%s", req.UserName, req.PassWord)
	session, err := service.Login(ctx, req)
	if err != nil {
		rsp.ResponseWithError(c, CodeLoginErr, err.Error())
		return
	}
	// 设置 cookie 值
	c.SetCookie(constant.SessionKey, session, constant.CookieExpire, "/", "", false, true)
	rsp.ResponseSuccess(c)
}

// Logout 登出
func Logout(c *gin.Context) {
	// 从上下文获取 session 会话 ID
	session, _ := c.Cookie(constant.SessionKey)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	req := &service.LogoutRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind get logout request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}

	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	// 带着 uuid 和 session 去操作redis 登出
	if err := service.Logout(ctx, req); err != nil {
		rsp.ResponseWithError(c, CodeLogoutErr, err.Error())
		return
	}

	// 将会话设置会过期，即浏览器删除该 cookie
	c.SetCookie(constant.SessionKey, session, -1, "/", "", false, true)
	rsp.ResponseSuccess(c)
}

// Logoff 注销
func Logoff(c *gin.Context) {
	req := &service.LogoffRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
	}
	// 获取 session 用于后续删除
	session, _ := c.Cookie(constant.SessionKey)
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)

	// 注销
	if err := service.Logoff(req, ctx); err != nil {
		log.Println("Logoff|Failed:%v", err)
		rsp.ResponseWithError(c, CodeLogoffErr, err.Error())
		return
	}

	// 注销成功，返回成功
	rsp.ResponseSuccess(c)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	//从HTTP请求中获取查询参数username和session
	userName := c.Query("username")
	session, _ := c.Cookie(constant.SessionKey)
	//使用context.WithValue函数创建一个带有session值的上下文，并将其存储在常量SessionKey中，以便在后续处理程序中使用。
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	req := &service.GetUserInfoRequest{
		UserName: userName,
	}
	fmt.Println(ctx)
	//创建一个HttpResponse结构体实例rsp用于返回响应
	rsp := &HttpResponse{}
	//生成一个uuid
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	//将其保存在请求上下文ctx
	ctx = context.WithValue(ctx, "uuid", uuid)
	fmt.Println(ctx)
	userInfo, err := service.GetUserInfo(ctx, req)
	if err != nil {
		rsp.ResponseWithError(c, CodeGetUserInfoErr, err.Error())
		return
	}
	rsp.ResponseWithData(c, userInfo)
}

// UpdateNickName 更新用户昵称
func UpdateNickName(c *gin.Context) {
	req := &service.UpdateNickNameRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind update user info request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	session, _ := c.Cookie(constant.SessionKey)
	log.Infof("UpdateNickName|session=%s", session)
	//使用context.WithValue函数创建一个带有session值的上下文，并将其存储在常量SessionKey中，以便在后续处理程序中使用。
	ctx := context.WithValue(context.Background(), constant.SessionKey, session)
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	if err := service.UpdateUserNickName(ctx, req); err != nil {
		rsp.ResponseWithError(c, CodeUpdateUserInfoErr, err.Error())
		return
	}
	rsp.ResponseSuccess(c)
}

// UpLoad 更新用户头像
func UpLoad(c *gin.Context) {

}
