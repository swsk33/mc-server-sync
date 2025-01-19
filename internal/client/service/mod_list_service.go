package service

import (
	"gitee.com/swsk33/mc-server-sync/internal/client/global"
	"gitee.com/swsk33/mc-server-sync/pkg/model"
	"gitee.com/swsk33/mc-server-sync/pkg/util"
	"gitee.com/swsk33/sclog"
	"net/http"
	"path/filepath"
)

// 用于计算模组列表的服务逻辑

// ModInfoMap 模组列表哈希表类型，其中：
//
//   - 键：模组文件的sha256摘要
//   - 值：完整模组文件元数据
type ModInfoMap map[string]*model.ModFile

// GetLocalModList 从本地获取模组列表
func GetLocalModList() (ModInfoMap, error) {
	// 获取模组文件夹下模组列表
	modList, e := model.NewModListFromFolder(global.TotalConfig.ModFolder)
	if e != nil {
		return nil, e
	}
	// 整理结果
	result := make(ModInfoMap)
	for _, item := range modList {
		result[item.Sha256] = item
	}
	return result, nil
}

// GetServerModList 从服务端获取模组列表
func GetServerModList() (ModInfoMap, error) {
	// 发送请求
	response, e := global.SendRequest("/api/mod-info/get-all", http.MethodGet, nil)
	if e != nil {
		sclog.ErrorLine("发送模组列表请求出错！")
		return nil, e
	}
	// 解析结果
	modList, e := model.ParseResultJson[[]*model.ModFile](response)
	if e != nil {
		sclog.ErrorLine("请求模组列表失败！")
		return nil, e
	}
	// 整理结果
	result := make(ModInfoMap)
	for _, item := range modList {
		result[item.Sha256] = item
	}
	return result, nil
}

// ExcludeModList 排除相关模组文件
//
//   - client 客户端本地模组文件列表
//   - server 服务端获取的模组文件列表
func ExcludeModList(client, server ModInfoMap) {
	// 根据配置获取排除的文件信息，并排除
	for _, name := range global.TotalConfig.Sync.IgnoreFileNames {
		path := filepath.Join(global.TotalConfig.ModFolder, name)
		if util.FileExists(path) {
			modInfo, e := model.NewModFile(path)
			if e != nil {
				sclog.ErrorLine(e.Error())
				continue
			}
			delete(client, modInfo.Sha256)
			delete(server, modInfo.Sha256)
		}
	}
}

// GetDownloadModList 计算需要从服务器下载的模组信息列表
// 即服务端存在，但是客户端不存在的模组列表
//
//   - client 客户端本地模组文件列表
//   - server 服务端获取的模组文件列表
//
// 返回需要下载的模组列表，模组文件名以服务端为基准
func GetDownloadModList(client, server ModInfoMap) []*model.ModFile {
	list := make([]*model.ModFile, 0)
	for checksum, info := range server {
		if _, ok := client[checksum]; !ok {
			list = append(list, info)
		}
	}
	return list
}

// GetRemovedModList 计算需要从本地移除的模组信息列表
// 即服务端不存在，但是客户端存在的模组列表
//
//   - client 客户端本地模组文件列表
//   - server 服务端获取的模组文件列表
//
// 返回需要移除的模组列表，模组文件名以客户端为基准
func GetRemovedModList(client, server ModInfoMap) []*model.ModFile {
	list := make([]*model.ModFile, 0)
	for checksum, info := range client {
		if _, ok := server[checksum]; !ok {
			list = append(list, info)
		}
	}
	return list
}