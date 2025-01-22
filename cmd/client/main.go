package main

import (
	"gitee.com/swsk33/gopher-fetch"
	"gitee.com/swsk33/mc-server-sync/internal/client/initialize"
	"gitee.com/swsk33/mc-server-sync/internal/client/service"
	"gitee.com/swsk33/sclog"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// 配置文件位置
var configPath string

// 启动服务端的子命令
var rootCmd = &cobra.Command{
	Use:   "mc-sync-client",
	Short: "Minecraft模组同步-客户端",
	Long:  "用于同步Minecraft模组的程序，该命令用于启动同步客户端",
	Run: func(cmd *cobra.Command, args []string) {
		startup()
	},
}

func init() {
	// 初始化下载器
	gopher_fetch.ConfigEnableLogger(false)
	gopher_fetch.ConfigEnvironmentProxy()
	// 初始化命令行
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "指定配置文件以启动客户端端")
}

func startup() {
	// 读取配置
	e := initialize.InitClientConfig(configPath)
	if e != nil {
		sclog.ErrorLine(e.Error())
		return
	}
	// 读取本地模组
	sclog.InfoLine("正在获取本地模组...")
	clientModMap, e := service.GetLocalModList()
	if e != nil {
		sclog.ErrorLine("获取本地模组失败！")
		sclog.ErrorLine(e.Error())
		return
	}
	// 获取服务端模组
	sclog.InfoLine("正在获取同步服务器模组列表...")
	serverModMap, e := service.GetServerModList()
	if e != nil {
		sclog.ErrorLine("获取服务器模组失败！")
		sclog.ErrorLine(e.Error())
		return
	}
	// 排除模组列表
	service.ExcludeModList(clientModMap, serverModMap)
	// 从服务端同步模组
	sclog.InfoLine("从服务端同步模组...")
	fetchList := service.GetDownloadModList(clientModMap, serverModMap)
	e = service.FetchModFromServer(fetchList)
	if e != nil {
		sclog.ErrorLine("从服务端同步模组失败！")
		sclog.ErrorLine(e.Error())
		return
	}
	// 移除本地多余的模组
	sclog.InfoLine("移除本地多余的模组...")
	removeList := service.GetRemovedModList(clientModMap, serverModMap)
	e = service.RemoveModFromLocal(removeList)
	if e != nil {
		sclog.ErrorLine("移除本地多余模组时发生错误！")
		sclog.ErrorLine(e.Error())
		return
	}
	sclog.InfoLine("同步工作已全部完成！")
}

func main() {
	// 视情况使用Cobra命令行逻辑或者直接启动
	if len(os.Args) < 2 {
		startup()
	} else {
		e := rootCmd.Execute()
		if e != nil {
			sclog.ErrorLine("执行客户端启动命令出错！")
			sclog.ErrorLine(e.Error())
		}
	}
	sclog.InfoLine("将在2s后退出...")
	time.Sleep(2 * time.Second)
}