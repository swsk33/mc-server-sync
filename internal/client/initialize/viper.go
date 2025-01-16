package initialize

import (
	"fmt"
	"gitee.com/swsk33/mc-server-sync/internal/client/global"
	"gitee.com/swsk33/mc-server-sync/pkg/model"
	"gitee.com/swsk33/sclog"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// InitClientConfig 初始化Viper及其配置对象
func InitClientConfig() error {
	// 从命令行参数获取配置文件路径
	// 如果未传递，则默认搜索当前路径下和可执行文件所在目录下的client-config.yaml文件
	if len(os.Args) < 2 {
		selfPath, e := os.Executable()
		if e != nil {
			return e
		}
		viper.SetConfigName("client-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Dir(selfPath))
	} else {
		viper.SetConfigFile(os.Args[1])
	}
	// 读取配置
	e := viper.ReadInConfig()
	if e != nil {
		sclog.ErrorLine("读取配置文件出错！")
		return e
	}
	// 反序列化配置
	e = viper.Unmarshal(&global.TotalConfig)
	if e != nil {
		sclog.ErrorLine("反序列化配置出错！")
		return e
	}
	// 设定默认值
	model.SetDefaultValue(&global.TotalConfig)
	if global.TotalConfig.ModFolder == "" {
		global.TotalConfig.ModFolder = ".minecraft/mods"
	}
	sclog.InfoLine("客户端已完成配置加载！")
	fmt.Println("客户端配置如下：")
	model.PrintConfig(global.TotalConfig, "")
	fmt.Println()
	return nil
}