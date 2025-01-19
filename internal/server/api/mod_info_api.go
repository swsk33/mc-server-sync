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

// List 列出全部模组文件信息
func (api *ModInfoApi) List(context *gin.Context) {
	// 获取模组文件夹下模组列表
	modList, e := model.NewModListFromFolder(global.TotalConfig.ModFolder)
	if e != nil {
		context.JSON(http.StatusOK, model.CreateFailedResult(fmt.Sprintf("读取模组文件夹失败！%s", e.Error())))
		return
	}
	// 返回
	context.JSON(http.StatusOK, model.CreateSuccessResult("获取模组文件列表成功！", modList))
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