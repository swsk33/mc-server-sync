package model

import (
	"fmt"
	"gitee.com/swsk33/mc-server-sync/pkg/util"
	"os"
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