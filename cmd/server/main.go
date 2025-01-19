package main

import (
	"gitee.com/swsk33/mc-server-sync/internal/server/initialize"
	"gitee.com/swsk33/sclog"
	"github.com/spf13/cobra"
	"os"
)

// 配置文件位置
var configPath string

// 启动服务端的子命令
var rootCmd = &cobra.Command{
	Use:   "mc-sync-server",
	Short: "Minecraft模组同步-服务端",
	Long:  "用于同步Minecraft模组的轻量级服务器，该命令用于启动同步服务端",
	Run: func(cmd *cobra.Command, args []string) {
		startup()
	},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "指定配置文件以启动服务端")
}

// 启动服务端的逻辑
func startup() {
	// 初始化配置
	e := initialize.InitServerConfig(configPath)
	if e != nil {
		sclog.ErrorLine(e.Error())
		return
	}
	// 启动服务
	e = initialize.InitGinRouterAndRun()
	if e != nil {
		sclog.ErrorLine("启动Gin路由出错！")
		sclog.ErrorLine(e.Error())
	}
}

func main() {
	// 视情况使用Cobra命令行逻辑或者直接启动
	if len(os.Args) < 2 {
		startup()
	} else {
		e := rootCmd.Execute()
		if e != nil {
			sclog.ErrorLine("执行服务端启动命令出错！")
			sclog.ErrorLine(e.Error())
		}
	}
}