package util

import "os"

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