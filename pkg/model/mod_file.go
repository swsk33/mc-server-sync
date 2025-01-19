package model

import (
	"fmt"
	"gitee.com/swsk33/mc-server-sync/pkg/util"
	"gitee.com/swsk33/sclog"
	"os"
	"path/filepath"
)

// ModFile 模组文件信息
type ModFile struct {
	// 文件名
	Filename string `json:"filename"`
	// 文件SHA256摘要值
	Sha256 string `json:"sha256"`
	// 文件大小（字节）
	Size int64 `json:"size"`
}

// NewModFile 指定模组文件完整路径，读取并返回其模组文件信息对象
//
//   - filepath 模组文件完整路径
//
// 返回模组文件信息
func NewModFile(filepath string) (*ModFile, error) {
	// 打开文件
	file, e := os.Open(filepath)
	if e != nil {
		return nil, e
	}
	defer func() {
		_ = file.Close()
	}()
	// 文件信息
	status, e := file.Stat()
	if e != nil {
		return nil, e
	}
	if status.IsDir() {
		return nil, fmt.Errorf("%s 是目录，不是模组文件！", filepath)
	}
	// 文件摘要
	checksum, e := util.CheckFileSHA256(filepath)
	if e != nil {
		return nil, e
	}
	// 返回文件信息
	return &ModFile{
		Filename: status.Name(),
		Sha256:   checksum,
		Size:     status.Size(),
	}, nil
}

// NewModListFromFolder 从文件夹读取其中全部文件作为模组列表
//
//   - modFolder 存放模组的文件夹
//
// 返回模组文件信息列表
func NewModListFromFolder(modFolder string) ([]*ModFile, error) {
	// 获取模组文件夹下文件
	fileList, e := os.ReadDir(modFolder)
	if e != nil {
		sclog.ErrorLine("读取本地模组列表出错！")
		return nil, e
	}
	// 获取文件信息
	modList := make([]*ModFile, 0)
	for _, file := range fileList {
		if !file.IsDir() {
			modFile, e := NewModFile(filepath.Join(modFolder, file.Name()))
			if e != nil {
				sclog.Error("读取模组文件：%s失败！原因：%s\n", file.Name(), e.Error())
				continue
			}
			modList = append(modList, modFile)
		}
	}
	return modList, nil
}