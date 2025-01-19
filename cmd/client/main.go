package main

import (
	"gitee.com/swsk33/mc-server-sync/internal/client/initialize"
	"gitee.com/swsk33/sclog"
	"github.com/spf13/cobra"
	"os"
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
}