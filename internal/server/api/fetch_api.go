package api

import (
	"fmt"
	"gitee.com/swsk33/mc-server-sync/internal/server/global"
	"gitee.com/swsk33/mc-server-sync/pkg/model"
	"gitee.com/swsk33/mc-server-sync/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"sync"
)

// FetchApi 下载文件的API对象
type FetchApi struct{}

// GetFileByName 根据文件名下载文件
func (api *FetchApi) GetFileByName(context *gin.Context) {
	// 获取文件
	name := context.Param("name")
	path := filepath.Join(global.TotalConfig.ModFolder, name)
	if !util.FileExists(path) {
		context.JSON(http.StatusNotFound, model.CreateFailedResult(fmt.Sprintf("找不到文件：%s", name)))
		return
	}
	// 返回文件
	context.File(path)
}

// API对象单例
var fetchApiInstance *FetchApi

// 单次初始化对象
var fetchApiOnce sync.Once

// GetFetchApiInstance 获取下载文件的API对象的全局唯一单例
func GetFetchApiInstance() *FetchApi {
	if fetchApiInstance == nil {
		fetchApiOnce.Do(func() {
			fetchApiInstance = &FetchApi{}
		})
	}
	return fetchApiInstance
}