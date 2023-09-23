package main

import (
	"Gous/config"
	"Gous/internal/router"
)

func Init() {
	config.InitConfig() // 初始化配置
}

func main() {
	Init()                       // 初始化信息
	router.InitRouterAndServer() // 路由配置、启动服务
}
