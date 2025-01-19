package api

import (
	"gitee.com/swsk33/mc-server-sync/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

// PingApi 用于测试连通性的Ping API
type PingApi struct{}

// Ping 测试连通性
func (api *PingApi) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, model.CreateNullSuccessResult("OK"))
}

// 唯一单例
var pingApiInstance *PingApi

// 单次初始化对象
var pingApiOnce sync.Once

// GetPingApiInstance 获取Ping Api唯一单例
func GetPingApiInstance() *PingApi {
	if pingApiInstance == nil {
		pingApiOnce.Do(func() {
			pingApiInstance = &PingApi{}
		})
	}
	return pingApiInstance
}