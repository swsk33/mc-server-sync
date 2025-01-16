package initialize

import (
	"fmt"
	"gitee.com/swsk33/mc-server-sync/internal/server/api"
	"gitee.com/swsk33/mc-server-sync/internal/server/global"
	"gitee.com/swsk33/sclog"
	"github.com/gin-gonic/gin"
)

// InitGinRouterAndRun 初始化Gin路由对象的函数
func InitGinRouterAndRun() error {
	gin.SetMode(gin.ReleaseMode)
	// 创建路由对象
	router := gin.Default()
	// 定义路由组
	fetchApi := api.GetFetchApiInstance()
	fetchGroup := router.Group("/api/fetch")
	{
		fetchGroup.GET("/get/:name", fetchApi.GetFileByName)
	}
	modInfoApi := api.GetModListApiInstance()
	modInfoGroup := router.Group("/api/mod-info")
	{
		modInfoGroup.GET("/get-all", modInfoApi.List)
	}
	// 启动服务
	sclog.InfoLine("启动Web服务器...")
	return router.Run(fmt.Sprintf(":%d", global.TotalConfig.Port))
}