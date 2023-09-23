package router

import (
	api "Gous/api/http/v1"
	"Gous/config"
	"Gous/pkg/constant"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// InitRouterAndServer 路由配置、启动服务
func InitRouterAndServer() {
	// 设置运行模式
	setAppRunMode()

	// 路由配置
	r := gin.Default()

	// 健康检查
	r.GET("/ping", api.Ping)
	// 用户注册
	r.POST("/user/register", api.Register)
	// 用户登录
	r.POST("/user/login", api.Login)
	// 用户登出
	r.POST("/user/logout", AuthMiddleWare(), api.Logout)
	// 用户注销
	r.POST("/user/logoff", AuthMiddleWare(), api.Logoff)
	// 获取用户信息
	r.GET("/user/get_user_info", AuthMiddleWare(), api.GetUserInfo)
	// 更新用户信息
	r.POST("/user/update_nick_name", AuthMiddleWare(), api.UpdateNickName)
	// 更新用户头像
	r.POST("/user/upload", api.UpLoad)

	// 渲染页面
	r.Static("/static/", "./web/static")
	r.Static("/upload/images/", "./web/upload/images")

	// 启动 server
	port := config.GetGlobalConf().AppConfig.Port
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		log.Error("start server err:" + err.Error())
	}
}

// 根据配置文件的设置来设置运行模式
func setAppRunMode() {
	if config.GetGlobalConf().AppConfig.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if config.GetGlobalConf().AppConfig.RunMode == "test" {
		gin.SetMode(gin.TestMode)
	} else if config.GetGlobalConf().AppConfig.RunMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		panic("setAppRunMode err, check app: run_mode")
	}
}

// AuthMiddleWare 检测用户释放处理登录状态
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//使用了 c.Cookie(constant.SessionKey) 方法来获取名为 constant.SessionKey 的 cookie 的值和错误信息
		if session, err := c.Cookie(constant.SessionKey); err == nil {
			if session != "" { //没有错误信息且session不为空，说明已经登录
				c.Next()
				return
			}
		}
		// 返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
		c.Abort() //终止后续处理程序函数的执行
		return
	}
}
