package service

import (
	"fmt"
	tp "gitee.com/swsk33/concurrent-task-pool/v2"
	"gitee.com/swsk33/gopher-fetch"
	"gitee.com/swsk33/mc-server-sync/internal/client/global"
	"gitee.com/swsk33/mc-server-sync/pkg/model"
	"gitee.com/swsk33/mc-server-sync/pkg/util"
	"gitee.com/swsk33/sclog"
	"os"
	"path/filepath"
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
	// 创建任务池准备下载
	taskPool := tp.NewSimpleTaskPool(global.TotalConfig.Sync.FetchConcurrency, list,
		func(task *model.ModFile, pool *tp.TaskPool[*model.ModFile]) {
			// 执行下载
			sclog.MutexInfo("正在下载模组文件：%s\n", task.Filename)
			url := fmt.Sprintf("http://%s:%d/api/fetch/get/%s", global.TotalConfig.Server.Host, global.TotalConfig.Server.Port, task.Filename)
			savePath := filepath.Join(global.TotalConfig.ModFolder, task.Filename)
			fetchTask := gopher_fetch.NewSimpleMonoGetTask(url, savePath)
			e := fetchTask.Run()
			if e != nil {
				sclog.MutexError("下载%s失败！原因：%s，将稍后重试...\n", task.Filename, e.Error())
				pool.Retry(task)
				return
			}
			// 校验sha256
			check, e := fetchTask.CheckFile(gopher_fetch.ChecksumSha256, task.Sha256)
			if e != nil {
				sclog.MutexError("校验%s失败！原因：%s，将稍后重试...\n", task.Filename, e.Error())
				pool.Retry(task)
				return
			}
			if !check {
				sclog.MutexError("%s的SHA256与预期不符！将重新下载...\n", task.Filename)
				pool.Retry(task)
				return
			}
			sclog.MutexInfo("同步模组%s完成！\n", task.Filename)
		})
	taskPool.Start()
	sclog.InfoLine("全部模组同步完成！")
	return nil
}

// RemoveModFromLocal 从本地移除模组文件
//
//   - list 需要移除的本地模组列表
func RemoveModFromLocal(list []*model.ModFile) error {
	if len(list) == 0 {
		sclog.InfoLine("无需从本地移除任何模组！")
		return nil
	}
	// 根据配置判断删除还是移动模组
	if global.TotalConfig.Sync.SoftRemove {
		backupFolder := "mod-backup"
		e := util.MkdirIfNotExists(backupFolder)
		if e != nil {
			sclog.ErrorLine("创建备份文件夹失败！")
			return e
		}
		for _, modFile := range list {
			e := os.Rename(filepath.Join(global.TotalConfig.ModFolder, modFile.Filename), filepath.Join(backupFolder, modFile.Filename))
			if e != nil {
				sclog.Error("移动文件%s失败！原因：%s\n", modFile.Filename, e.Error())
				continue
			}
			sclog.Info("已移动模组文件%s到备份目录！\n", modFile.Filename)
		}
	} else {
		for _, modFile := range list {
			e := os.Remove(filepath.Join(global.TotalConfig.ModFolder, modFile.Filename))
			if e != nil {
				sclog.Error("删除文件%s失败！原因：%s\n", modFile.Filename, e.Error())
				continue
			}
			sclog.Info("已删除模组文件：%s\n", modFile.Filename)
		}
	}
	sclog.InfoLine("已移除全部多余模组！")
	return nil
}