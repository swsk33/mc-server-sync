package main

import (
	"gitee.com/swsk33/mc-server-sync/internal/server/initialize"
	"gitee.com/swsk33/sclog"
)

func main() {
	// 初始化配置
	e := initialize.InitServerConfig()
	if e != nil {
		sclog.ErrorLine(e.Error())
		return
	}
	// 启动服务
	e = initialize.InitGinRouterAndRun()
	if e != nil {
		sclog.ErrorLine("启动Gin路由出错！")
		sclog.ErrorLine(e.Error())
	}
}