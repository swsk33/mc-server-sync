package main

import (
	"gitee.com/swsk33/gopher-fetch"
	"gitee.com/swsk33/mc-server-sync/internal/client/global"
	"gitee.com/swsk33/mc-server-sync/internal/client/initialize"
	"gitee.com/swsk33/mc-server-sync/internal/client/service"
	"gitee.com/swsk33/mc-server-sync/pkg/param"
	"gitee.com/swsk33/mc-server-sync/pkg/util"
	"gitee.com/swsk33/sclog"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

var (
	// 配置文件位置
	configPath string
	// 强制运行目录为应用程序所在目录
	forceWorkDirectory bool
	// 是否在终端模拟器中显式执行
	inTerminal bool
)

// 启动服务端的子命令
var rootCmd = &cobra.Command{
	Use:   "mc-sync-client",
	Short: "Minecraft模组同步-客户端",
	Long:  "用于同步Minecraft模组的程序，该命令用于启动同步客户端",
	Run: func(cmd *cobra.Command, args []string) {
		// 执行启动逻辑
		e := startup()
		// 错误处理
		handleClientLaunchError(e)
	},
}

func init() {
	// 初始化下载器
	gopher_fetch.ConfigEnableLogger(false)
	gopher_fetch.ConfigEnvironmentProxy()
	// 初始化命令行
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "指定配置文件以启动客户端")
	rootCmd.Flags().BoolVarP(&forceWorkDirectory, "force-work-directory", "d", false, "若带上该标志，则会强制程序的工作目录为程序自身的所在目录")
	rootCmd.Flags().BoolVarP(&inTerminal, "in-terminal", "t", false, "若带上该标志，则会调用系统可用的终端模拟器程序（例如cmd、gnome-terminal等）弹出新的窗口运行客户端程序，建议使用游戏启动器调用同步客户端时加上该标志，使得同步过程以及日志能够显现")
}

// 调用终端显示运行客户端逻辑
func runInTerminal() error {
	args := util.RemoveArgs(os.Args, "-t", "--in-terminal")
	e := util.ExecuteByTerminal(args...)
	if e != nil {
		return e
	}
	return nil
}

func startup() error {
	// 若在终端运行，则重新调用终端执行客户端命令
	if inTerminal {
		e := runInTerminal()
		if e != nil {
			sclog.ErrorLine("唤起终端运行出现错误！")
		}
		return e
	}
	// 工作目录处理
	if forceWorkDirectory {
		selfPath, e := os.Executable()
		if e != nil {
			sclog.ErrorLine("获取自身路径失败！")
			return e
		}
		workDir := filepath.Dir(selfPath)
		e = os.Chdir(workDir)
		if e != nil {
			sclog.ErrorLine("修改工作目录失败！")
			return e
		}
		sclog.Warn("已改变程序工作目录为其自身所在目录：%s\n", workDir)
		sclog.WarnLine("请注意配置文件以及模组文件夹的相对位置！")
	}
	// 读取配置
	e := initialize.InitClientConfig(configPath)
	if e != nil {
		return e
	}
	// 读取本地模组
	sclog.InfoLine("正在获取本地模组...")
	clientModMap, e := service.GetLocalModList()
	if e != nil {
		sclog.ErrorLine("获取本地模组失败！")
		return e
	}
	// 获取服务端模组
	sclog.InfoLine("正在获取同步服务器模组列表...")
	serverModMap, e := service.GetServerModList()
	if e != nil {
		sclog.ErrorLine("获取服务器模组失败！")
		return e
	}
	// 排除模组列表
	service.ExcludeModList(clientModMap, serverModMap)
	// 从服务端同步模组
	sclog.InfoLine("从服务端同步模组...")
	fetchList := service.GetDownloadModList(clientModMap, serverModMap)
	e = service.FetchModFromServer(fetchList)
	if e != nil {
		sclog.ErrorLine("从服务端同步模组失败！")
		return e
	}
	// 移除本地多余的模组
	sclog.InfoLine("移除本地多余的模组...")
	removeList := service.GetRemovedModList(clientModMap, serverModMap)
	service.RemoveModFromLocal(removeList)
	// 最后，检查本地重复的模组文件并移除
	sclog.InfoLine("移除本地重复模组...")
	e = service.RemoveDuplicateModFromLocal()
	if e != nil {
		sclog.ErrorLine("移除本地多余模组时发生错误！")
		return e
	}
	sclog.InfoLine("同步工作已全部完成！")
	return nil
}

// 处理客户端启动错误，若e不为nil，则退出程序
func handleClientLaunchError(e error) {
	if e != nil {
		sclog.ErrorLine("启动同步客户端失败！")
		sclog.ErrorLine(e.Error())
		util.ErrorExitAndDelay(3)
	}
}

func main() {
	sclog.Info("模组同步-客户端 v%s，启动！\n", param.ClientVersion)
	// 错误对象
	var e error
	// 视情况使用Cobra命令行逻辑或者直接启动
	if len(os.Args) < 2 {
		e = startup()
	} else {
		e = rootCmd.Execute()
	}
	// 错误处理
	handleClientLaunchError(e)
	// 正常情况将按照配置文件延迟退出
	if global.TotalConfig.ExitDelay > 0 {
		sclog.Info("将在%ds后退出...\n", global.TotalConfig.ExitDelay)
		time.Sleep(time.Duration(global.TotalConfig.ExitDelay) * time.Second)
	}
}