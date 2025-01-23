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
	// 基本通用配置
	Base BaseConfig `mapstructure:"base"`
	// 服务端口
	Port int `mapstructure:"port" default:"25566"`
	// 守护进程模式运行时的配置
	Daemon struct {
		// 日志文件
		LogFile string `mapstructure:"log-file" default:"mc-sync-server.log"`
	} `mapstructure:"daemon"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	// 基本通用配置
	Base BaseConfig `mapstructure:"base"`
	// 连接同步服务器配置
	Server struct {
		// 服务器地址
		Host string `mapstructure:"host"`
		// 服务器端口
		Port int `mapstructure:"port" default:"25566"`
	} `mapstructure:"server"`
	// 文件同步相关配置
	Sync struct {
		// 同时下载的文件数
		FetchConcurrency int `mapstructure:"fetch-concurrency" default:"3"`
		// 软删除
		// 若开启软删除，则在同步时删除本地模组文件时，不会实际对文件进行删除，而是移动到程序所在目录的mod-backup文件夹
		SoftRemove bool `mapstructure:"soft-remove" default:"true"`
		// 忽略同步的模组文件名列表
		// 默认情况下，若本地存在但服务器不存在的模组会被删除
		// 若某文件被加入到了忽略列表，则即使服务器不存在该文件也不会被删除
		IgnoreFileNames []string `mapstructure:"ignore-file-names"`
	} `mapstructure:"sync"`
	// 退出延迟，同步完成后延迟多少秒退出，若设为0则同步完成立即退出
	// 建议设定延迟几秒，以便于查看同步日志，排查错误
	ExitDelay int `mapstructure:"exit-delay" default:"2"`
}

// SetDefaultValue 根据对象的default标签设定字段默认值
func SetDefaultValue(object any) {
	// 获取该对象的反射对象
	value := reflect.ValueOf(object)
	// 如果传入的是指针，获取指针指向的值
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	} else {
		// 如果传入的是值类型，强制转换为指针类型
		value = reflect.New(value.Type()).Elem()
	}
	// 获取该对象类型
	objectType := reflect.TypeOf(object)
	// 如果对象本身是指针类型，获取指针指向的类型
	if objectType.Kind() == reflect.Ptr {
		objectType = objectType.Elem()
	}
	// 获取属性个数
	fieldNum := value.NumField()
	// 遍历每个属性信息
	for i := 0; i < fieldNum; i++ {
		fieldType := objectType.Field(i)
		fieldValue := value.Field(i)
		// 结构体类型则递归
		if fieldType.Type.Kind() == reflect.Struct {
			SetDefaultValue(fieldValue.Addr().Interface())
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
				case reflect.Bool:
					if boolValue, e := strconv.ParseBool(tagValue); e == nil {
						fieldValue.SetBool(boolValue)
					}
				default:
					panic("不支持的类型")
				}
			}
		}
	}
	// 如果原始传入对象是值类型，需要将修改的值写回去
	if reflect.ValueOf(object).Kind() != reflect.Ptr {
		reflect.ValueOf(object).Elem().Set(value)
	}
}

// PrintConfig 输出配置信息到控制台
//
//   - config 打印的配置对象
//   - prefix 递归前缀
func PrintConfig(config any, prefix string) {
	// 获取该对象的反射对象
	value := reflect.ValueOf(config)
	// 获取该对象类型
	objectType := reflect.TypeOf(config)
	// 获取属性个数
	fieldNum := value.NumField()
	// 遍历每个属性信息
	for i := 0; i < fieldNum; i++ {
		// 获取字段的反射类型
		field := objectType.Field(i)
		// 获取字段值
		fieldValue := value.Field(i)
		// 判断是否是匿名字段
		var configName string
		if field.Anonymous {
			// 对于匿名字段，递归输出
			PrintConfig(fieldValue.Interface(), prefix)
			continue
		} else {
			// 如果字段有 `mapstructure` 标签，则使用标签的值作为配置项名称
			configName = field.Tag.Get("mapstructure")
			if configName == "" {
				// 如果没有 `mapstructure` 标签，则使用字段的名称
				configName = field.Name
			}
		}
		// 如果是结构体类型则递归
		if fieldValue.Kind() == reflect.Struct {
			if prefix == "" {
				PrintConfig(fieldValue.Interface(), configName)
			} else {
				PrintConfig(fieldValue.Interface(), fmt.Sprintf("%s.%s", prefix, configName))
			}
		} else {
			// 打印字段值
			if prefix == "" {
				fmt.Printf("%s = %v\n", configName, fieldValue.Interface())
			} else {
				fmt.Printf("%s = %v\n", fmt.Sprintf("%s.%s", prefix, configName), fieldValue.Interface())
			}
		}
	}
}