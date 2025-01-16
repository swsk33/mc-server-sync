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