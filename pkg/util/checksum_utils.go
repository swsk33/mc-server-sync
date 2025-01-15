package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

// CheckFileSHA256 计算文件的SHA256的值
//
//   - file 传入文件路径
//
// 返回文件的SHA256摘要
func CheckFileSHA256(filepath string) (string, error) {
	// 打开文件
	file, e := os.Open(filepath)
	if e != nil {
		return "", e
	}
	defer func() {
		_ = file.Close()
	}()
	// 校验SHA256
	hashChecker := sha256.New()
	// 计算摘要
	_, e = io.Copy(hashChecker, file)
	if e != nil {
		return "", e
	}
	return strings.ToLower(fmt.Sprintf("%x", hashChecker.Sum(nil))), e
}