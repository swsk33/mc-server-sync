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
//
//   - config 指定配置文件路径，若指定为空字符串""，则默认搜索当前路径下和可执行文件所在目录下的client-config.yaml文件
func InitClientConfig(config string) error {
	// 从命令行参数获取配置文件路径
	if config == "" {
		selfPath, e := os.Executable()
		if e != nil {
			return e
		}
		viper.SetConfigName("client-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Dir(selfPath))
	} else {
		viper.SetConfigFile(config)
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
	if global.TotalConfig.Base.ModFolder == "" {
		global.TotalConfig.Base.ModFolder = ".minecraft/mods"
	}
	sclog.InfoLine("客户端已完成配置加载！")
	fmt.Println("客户端配置如下：")
	model.PrintConfig(global.TotalConfig, "")
	fmt.Println()
	return nil
}