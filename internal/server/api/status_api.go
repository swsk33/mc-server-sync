package api

import (
	"gitee.com/swsk33/mc-server-sync/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sync"
)

// StatusApi 用于获取系统状态的API
type StatusApi struct{}

// Ping 测试连通性
func (api *StatusApi) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, model.CreateNullSuccessResult("OK"))
}

// DaemonPid 获取守护进程PID
func (api *StatusApi) DaemonPid(context *gin.Context) {
	context.JSON(http.StatusOK, model.CreateSuccessResult[int]("获取进程PID成功！", os.Getpid()))
}

// 唯一单例
var statusApiInstance *StatusApi

// 单次初始化对象
var statusApiOnce sync.Once

// GetStatusApiInstance 获取Ping Api唯一单例
func GetStatusApiInstance() *StatusApi {
	if statusApiInstance == nil {
		statusApiOnce.Do(func() {
			statusApiInstance = &StatusApi{}
		})
	}
	return statusApiInstance
}