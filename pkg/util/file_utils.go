package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists 判断文件是否存在
func FileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	if e == nil {
		return true
	}
	return !os.IsNotExist(e)
}

// MkdirIfNotExists 创建文件夹如果文件夹不存在
//
//   - folder 文件夹路径
func MkdirIfNotExists(folder string) error {
	if !FileExists(folder) {
		return os.MkdirAll(folder, 0755)
	}
	return nil
}

// MoveFile 移动文件，如果目标目录已存在同名文件，则自动重命名
//
//   - origin 原始文件路径
//   - destination 移动后文件路径
func MoveFile(origin, destination string) error {
	// 获取目标目录和文件名
	destFolder := filepath.Dir(destination)
	destFileName := filepath.Base(destination)
	// 截取目标文件的主要文件名部分
	extName := filepath.Ext(destFileName)
	baseName := destFileName[:len(destFileName)-len(extName)]
	// 检查目标文件是否存在，如果存在则重命名
	counter := 1
	newDestination := destination
	for {
		if !FileExists(newDestination) {
			break
		}
		newFileName := fmt.Sprintf("%s-%d%s", baseName, counter, extName)
		newDestination = filepath.Join(destFolder, newFileName)
		counter++
	}
	// 移动文件
	e := os.Rename(origin, newDestination)
	if e != nil {
		return fmt.Errorf("移动文件失败: %w", e)
	}
	return nil
}