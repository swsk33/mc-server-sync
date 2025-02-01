package main

import (
	"fmt"
	"gitee.com/swsk33/mc-server-sync/internal/server/global"
	"gitee.com/swsk33/mc-server-sync/internal/server/initialize"
	"gitee.com/swsk33/mc-server-sync/pkg/param"
	"gitee.com/swsk33/mc-server-sync/pkg/util"
	"gitee.com/swsk33/sclog"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"os"
)

var (
	// 配置文件位置
	configPath string
	// 是否以守护进程形式后台启动
	isDaemon bool
)

// 启动服务端的子命令
var rootCmd = &cobra.Command{
	Use:   "mc-sync-server",
	Short: "Minecraft模组同步-服务端",
	Long:  "用于同步Minecraft模组的轻量级服务器，该命令用于启动同步服务端",
	Run: func(cmd *cobra.Command, args []string) {
		// 执行启动逻辑
		e := startup()
		// 错误处理
		handleServerLaunchError(e)
	},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "指定配置文件以启动服务端")
	rootCmd.Flags().BoolVarP(&isDaemon, "daemon", "d", false, "若加上该标志，则服务端将以守护进程形式后台运行，不支持Windows操作系统")
}

// 在守护进程启动Web服务器后台运行
func startDaemon() error {
	// 处理命令行，移除-d和--daemon
	args := util.RemoveArgs(os.Args, "-d", "--daemon")
	// 创建守护进程对象，包含参数等
	daemonContext := &daemon.Context{
		// 日志文件路径
		LogFileName: global.TotalConfig.Daemon.LogFile,
		// 日志文件权限
		LogFilePerm: 0640,
		// 工作目录
		WorkDir: ".",
		// 文件掩码
		Umask: 027,
		// 启动运行参数
		Args: args,
	}
	// 创建守护进程
	// Reborn函数将根据上述给定运行参数Args，重新运行一个进程并以守护进程形式后台运行
	// 返回一个Process对象和错误对象
	daemonProcess, e := daemonContext.Reborn()
	if e != nil {
		sclog.ErrorLine("启动守护进程出错！")
		return e
	}
	// Reborn方法返回的Process对象是其创建的守护进程信息，可以打印
	// 到此，主进程可以结束了，守护进程将在后台一直运行
	if daemonProcess != nil {
		sclog.InfoLine("服务端已在守护进程中运行！")
		sclog.Info("PID：%d 日志文件：%s\n", daemonProcess.Pid, global.TotalConfig.Daemon.LogFile)
	}
	return nil
}

// 启动服务端的逻辑
func startup() error {
	// 初始化配置
	e := initialize.InitServerConfig(configPath)
	if e != nil {
		return e
	}
	// 检查目录
	if !util.FileExists(global.TotalConfig.Base.ModFolder) {
		sclog.Error("模组文件夹：%s不存在！请配置正确的模组文件夹！\n", global.TotalConfig.Base.ModFolder)
		return fmt.Errorf("模组文件夹：%s不存在！\n", global.TotalConfig.Base.ModFolder)
	}
	if !util.FileExists(global.TotalConfig.ClientModFolder) {
		sclog.Warn("仅客户端类型模组文件夹：%s不存在！将创建该文件夹...\n", global.TotalConfig.ClientModFolder)
		e := util.MkdirIfNotExists(global.TotalConfig.ClientModFolder)
		if e != nil {
			sclog.Error("创建文件夹：%s失败！\n", global.TotalConfig.ClientModFolder)
			return e
		}
	}
	// 启动服务
	if isDaemon {
		e = startDaemon()
	} else {
		e = initialize.InitGinRouterAndRun()
	}
	if e != nil {
		sclog.ErrorLine("启动同步Web服务器出错！")
		return e
	}
	return nil
}

// 处理服务端启动错误，若e不为nil，则退出程序
func handleServerLaunchError(e error) {
	if e != nil {
		sclog.ErrorLine("启动同步服务端失败！")
		sclog.ErrorLine(e.Error())
		util.ErrorExitAndDelay(3)
	}
}

func main() {
	sclog.Info("模组同步-服务端 v%s，启动！\n", param.ServerVersion)
	// 错误对象
	var e error
	// 视情况使用Cobra命令行逻辑或者直接启动
	if len(os.Args) < 2 {
		e = startup()
	} else {
		e = rootCmd.Execute()
	}
	// 错误处理
	handleServerLaunchError(e)
}