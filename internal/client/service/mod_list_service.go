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

// GetLocalModList 从本地获取模组列表
//
// 分别返回：
//   - 本地模组列表（已去重）
//   - 重复的模组文件列表
//   - 错误对象
func GetLocalModList() (model.ModInfoMap, []*model.ModFile, error) {
	// 获取模组文件夹下模组列表
	modList, e := model.NewModListFromFolder(global.TotalConfig.Base.ModFolder, model.ModTypeBoth)
	if e != nil {
		return nil, nil, e
	}
	// 返回结果
	modMap, duplicates := model.NewModInfoMapFromSlice(modList)
	return modMap, duplicates, nil
}

// GetServerModList 从服务端获取需要同步的模组列表
//
// 分别返回：
//   - 服务端模组列表
//   - 错误对象
func GetServerModList() (model.ModInfoMap, error) {
	// 发送请求
	response, e := global.SendRequest("/api/mod-info/get-all-both", http.MethodGet, nil)
	if e != nil {
		sclog.ErrorLine("发送双端模组列表请求出错！")
		return nil, e
	}
	// 解析结果
	modList, e := model.ParseResultJson[[]*model.ModFile](response)
	if e != nil {
		sclog.ErrorLine("解析双端模组列表请求失败！")
		return nil, e
	}
	// 根据配置，获取服务端提供的仅客户端类型模组列表并合并到服务端模组列表
	if global.TotalConfig.Sync.FetchClientMods {
		// 发送请求
		clientModResponse, e := global.SendRequest("/api/mod-info/get-all-client", http.MethodGet, nil)
		if e != nil {
			sclog.ErrorLine("发送仅客户端模组列表请求出错！")
			return nil, e
		}
		// 解析结果
		clientModList, e := model.ParseResultJson[[]*model.ModFile](clientModResponse)
		if e != nil {
			sclog.ErrorLine("解析仅客户端模组列表请求失败！")
			return nil, e
		}
		// 合并结果
		modList = append(modList, clientModList...)
	} else {
		sclog.WarnLine("未开启仅客户端类型模组同步，将不会从服务端获取并同步服务器提供的客户端类型模组列表！")
	}
	// 返回结果
	modMap, _ := model.NewModInfoMapFromSlice(modList)
	return modMap, nil
}

// ExcludeModList 排除相关模组文件
//
//   - client 客户端本地模组文件列表
//   - server 服务端获取的模组文件列表
func ExcludeModList(client, server model.ModInfoMap) {
	// 根据配置获取排除的文件信息，并排除
	for _, name := range global.TotalConfig.Sync.IgnoreFileNames {
		path := filepath.Join(global.TotalConfig.Base.ModFolder, name)
		if util.FileExists(path) {
			modInfo, e := model.NewModFile(path, model.ModTypeBoth)
			if e != nil {
				sclog.ErrorLine(e.Error())
				continue
			}
			client.Remove(modInfo)
			server.Remove(modInfo)
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
func GetDownloadModList(client, server model.ModInfoMap) []*model.ModFile {
	return server.Subtract(client)
}

// GetRemovedModList 计算需要从本地移除的模组信息列表
// 即服务端不存在，但是客户端存在的模组列表
//
//   - client 客户端本地模组文件列表
//   - server 服务端获取的模组文件列表
//
// 返回需要移除的模组列表，模组文件名以客户端为基准
func GetRemovedModList(client, server model.ModInfoMap) []*model.ModFile {
	return client.Subtract(server)
}