package main

import (
	"gitee.com/swsk33/mc-server-sync/internal/client/initialize"
	"gitee.com/swsk33/sclog"
)

func main() {
	// 读取配置
	e := initialize.InitClientConfig()
	if e != nil {
		sclog.ErrorLine(e.Error())
		return
	}
}