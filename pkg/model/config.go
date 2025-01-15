package model

import (
	"fmt"
	"reflect"
	"strconv"
)

// BaseConfig 客户端与服务端的通用配置模型
type BaseConfig struct {
	// 模组文件夹位置
	// 服务端默认模组文件夹是当前运行路径下的mods
	// 客户端默认模组文件夹是当前运行路径下的.minecraft/mods
	ModFolder string `mapstructure:"mod-folder"`
}

// ServerConfig 服务端配置
type ServerConfig struct {
	BaseConfig
	// 服务端口
	Port int `mapstructure:"port" default:"25566"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	BaseConfig
	// 连接同步服务器配置
	SyncServer struct {
		// 服务器地址
		Host string `mapstructure:"host"`
		// 服务器端口
		Port int `mapstructure:"port" default:"25566"`
	} `mapstructure:"sync-server"`
	// 下载相关配置
	Download struct {
		// 同时下载的文件数
		Concurrency int `mapstructure:"concurrency" default:"3"`
	} `mapstructure:"download"`
	// 忽略同步的模组文件名列表
	IgnoreFileNames []string `mapstructure:"ignore-file-names"`
}

// SetDefaultValue 根据对象的default标签设定字段默认值
func SetDefaultValue(object any) {
	// 获取该对象的反射对象
	value := reflect.ValueOf(object)
	// 获取该对象类型
	objectType := reflect.TypeOf(object)
	// 获取属性个数
	fieldNum := value.NumField()
	// 遍历每个属性信息
	for i := 0; i < fieldNum; i++ {
		fieldType := objectType.Field(i)
		fieldValue := value.Field(i)
		// 结构体类型则递归
		if fieldType.Type.Kind() == reflect.Struct {
			SetDefaultValue(fieldValue.Interface())
		} else if fieldValue.IsZero() {
			// 字段为空则赋值默认值
			tagValue := fieldType.Tag.Get("default")
			if tagValue != "" {
				switch fieldType.Type.Kind() {
				case reflect.String:
					fieldValue.SetString(tagValue)
				case reflect.Int:
					if intValue, e := strconv.Atoi(tagValue); e == nil {
						fieldValue.SetInt(int64(intValue))
					}
				default:
					panic("不支持的类型")
				}
			}
		}
	}
}

// PrintConfig 输出配置信息到控制台
func PrintConfig(config any) {
	fmt.Println("配置信息如下：")
	// 获取该对象的反射对象
	value := reflect.ValueOf(config)
	// 获取该对象类型
	objectType := reflect.TypeOf(config)
	// 获取属性个数
	fieldNum := value.NumField()
	// 遍历每个属性信息
	for i := 0; i < fieldNum; i++ {
		configName := objectType.Field(i).Tag.Get("mapstructure")
		configValue := value.Field(i).Interface()
		fmt.Printf("%s = %s\n", configName, configValue)
	}
}