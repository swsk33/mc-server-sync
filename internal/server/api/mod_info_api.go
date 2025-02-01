package api

import (
	"fmt"
	"gitee.com/swsk33/mc-server-sync/internal/server/global"
	"gitee.com/swsk33/mc-server-sync/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

// ModInfoApi 用于获取模组文件信息的API对象
type ModInfoApi struct{}

// ListBoth 列出全部双端模组文件信息
// 该API用于获取安装在Minecraft服务端的模组，这些模组通常服务端和客户端都需要安装
func (api *ModInfoApi) ListBoth(context *gin.Context) {
	// 获取模组文件夹下模组列表
	modList, e := model.NewModListFromFolder(global.TotalConfig.Base.ModFolder, model.ModTypeBoth)
	if e != nil {
		context.JSON(http.StatusOK, model.CreateFailedResult(fmt.Sprintf("读取模组文件夹失败！%s", e.Error())))
		return
	}
	// 返回
	context.JSON(http.StatusOK, model.CreateSuccessResult("获取模组文件列表成功！", modList))
}

// ListClient 列出全部客户端类型模组文件信息
// 该API用于获取单独配置提供的仅客户端模组
func (api *ModInfoApi) ListClient(context *gin.Context) {
	// 获取模组文件夹下模组列表
	modList, e := model.NewModListFromFolder(global.TotalConfig.ClientModFolder, model.ModTypeClient)
	if e != nil {
		context.JSON(http.StatusOK, model.CreateFailedResult(fmt.Sprintf("读取客户端模组文件夹失败！%s", e.Error())))
		return
	}
	// 返回
	context.JSON(http.StatusOK, model.CreateSuccessResult("获取客户端模组文件列表成功！", modList))
}

// API对象单例
var modInfoApiInstance *ModInfoApi

// 单次初始化对象
var modInfoApiOnce sync.Once

// GetModListApiInstance 获取下载文件的API对象的全局唯一单例
func GetModListApiInstance() *ModInfoApi {
	if modInfoApiInstance == nil {
		modInfoApiOnce.Do(func() {
			modInfoApiInstance = &ModInfoApi{}
		})
	}
	return modInfoApiInstance
}