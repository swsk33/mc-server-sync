package service

import (
	"fmt"
	tp "gitee.com/swsk33/concurrent-task-pool/v2"
	"gitee.com/swsk33/mc-server-sync/internal/client/global"
	"gitee.com/swsk33/mc-server-sync/pkg/model"
	"gitee.com/swsk33/sclog"
)

// 同步模组文件的服务逻辑

// FetchModFromServer 从服务器下载模组到本地
//
//   - list 要下载的模组列表
func FetchModFromServer(list []*model.ModFile) error {
	if len(list) == 0 {
		sclog.InfoLine("无需从服务端同步下载模组！")
		return nil
	}
	// 计算下载地址
	fetchPathList := make([]string, 0)
	for _, mod := range list {
		fetchPathList = append(fetchPathList, fmt.Sprintf("http://%s:%d/api/fetch/get/%s", global.TotalConfig.Server.Host, global.TotalConfig.Server.Port, mod.Filename))
	}
	// 创建任务池准备下载
	taskPool := tp.NewSimpleTaskPool(global.TotalConfig.Sync.FetchConcurrency, fetchPathList,
		func(task string, pool *tp.TaskPool[string]) {
			// 创建下载任务

		})
	taskPool.Start()
	return nil
}