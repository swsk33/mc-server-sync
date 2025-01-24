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

// 根据配置，移除或者移动某个本地模组文件
//
//   - mod 要移除的模组文件对象
func removeModFile(mod *model.ModFile) error {
	// 根据配置判断删除还是移动模组
	if global.TotalConfig.Sync.SoftRemove {
		// 准备模组备份文件夹
		backupFolder := "mod-backup"
		e := util.MkdirIfNotExists(backupFolder)
		if e != nil {
			sclog.ErrorLine("创建备份文件夹失败！")
			return e
		}
		// 移动模组文件
		e = os.Rename(filepath.Join(global.TotalConfig.Base.ModFolder, mod.Filename), filepath.Join(backupFolder, mod.Filename))
		if e != nil {
			sclog.Error("移动文件%s失败！\n", mod.Filename)
			return e
		}
		sclog.Info("已移动模组文件%s到备份目录！\n", mod.Filename)
	} else {
		// 删除文件
		e := os.Remove(filepath.Join(global.TotalConfig.Base.ModFolder, mod.Filename))
		if e != nil {
			sclog.Error("删除文件%s失败！\n", mod.Filename)
			return e
		}
		sclog.Info("已删除模组文件：%s\n", mod.Filename)
	}
	return nil
}

// FetchModFromServer 从服务器下载缺少的模组到本地
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
			savePath := filepath.Join(global.TotalConfig.Base.ModFolder, task.Filename)
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

// RemoveModFromLocal 从本地移除多余模组文件
//
//   - list 需要移除的本地模组列表
func RemoveModFromLocal(list []*model.ModFile) {
	if len(list) == 0 {
		sclog.InfoLine("无需从本地移除任何模组！")
		return
	}
	// 执行移除操作
	for _, modFile := range list {
		e := removeModFile(modFile)
		if e != nil {
			sclog.ErrorLine(e.Error())
		}
	}
	sclog.InfoLine("已移除全部多余模组！")
	return
}

// RemoveDuplicateModFromLocal 从本地移除重复的模组
// 将会移除本地文件名不同、但是SHA256相同的模组文件
func RemoveDuplicateModFromLocal() error {
	// 再次获取本地模组列表
	list, e := model.NewModListFromFolder(global.TotalConfig.Base.ModFolder)
	if e != nil {
		sclog.ErrorLine("获取本地模组列表出错！")
		return e
	}
	// 查找出重复的
	duplicates := model.FindDuplicate(list)
	if len(duplicates) == 0 {
		sclog.InfoLine("本地没有重复的模组，无需移除！")
		return nil
	}
	// 进行移除操作
	for _, item := range duplicates {
		sclog.Warn("发现本地有重复模组文件：%s，将进行移除...\n", item.Filename)
		e := removeModFile(item)
		if e != nil {
			sclog.ErrorLine(e.Error())
		}
	}
	sclog.InfoLine("已移除全部本地重复的模组文件！")
	return nil
}