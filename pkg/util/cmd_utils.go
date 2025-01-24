package util

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// RemoveArgs 从命令行参数切片中移除某些参数
//
//   - origin 原始参数列表
//   - remove 要移除的参数列表
//
// 返回移除对应参数后的参数列表
func RemoveArgs(origin []string, remove ...string) []string {
	// 转换为Map
	originMap, removeMap := make(map[string]any), make(map[string]any)
	for _, item := range origin {
		originMap[item] = nil
	}
	for _, item := range remove {
		removeMap[item] = nil
	}
	// 进行移除
	result := make([]string, 0)
	for arg := range originMap {
		if _, ok := removeMap[arg]; ok {
			continue
		}
		result = append(result, arg)
	}
	return result
}

// CheckCommand 检查某个命令是否存在
//
//   - exe 检查的命令
//
// 若命令存在返回true
func CheckCommand(exe string) bool {
	_, e := exec.LookPath(exe)
	return e == nil
}

// GetSystemTerminal 检查系统的终端
//
// 返回系统安装的一个终端程序，否则返回错误对象
func GetSystemTerminal() (string, error) {
	linuxTerminals := []string{"gnome-terminal", "konsole", "deepin-terminal", "xterm", "kgx", "tilix"}
	for _, terminal := range linuxTerminals {
		if CheckCommand(terminal) {
			return terminal, nil
		}
	}
	return "", fmt.Errorf("没有找到支持的终端程序！请安装下列终端模拟器的其中之一：%s", strings.Join(linuxTerminals, ", "))
}

// ExecuteByTerminal 调用系统默认终端执行一行命令
// 通过这种方式可以实现弹出终端执行对应命令，使命令的执行过程能够显现
//
//   - args 运行的命令行
//
// 出现错误则返回错误对象
func ExecuteByTerminal(args ...string) error {
	// 对于Windows操作系统
	if runtime.GOOS == "windows" {
		totalArgs := []string{"/c", "start", "/wait"}
		totalArgs = append(totalArgs, args...)
		cmd := exec.Command("cmd", totalArgs...)
		return cmd.Run()
	}
	// 对于Linux操作系统
	if runtime.GOOS == "linux" {
		terminal, e := GetSystemTerminal()
		if e != nil {
			return e
		}
		// 组装参数
		totalArgs := make([]string, 0)
		switch terminal {
		case "gnome-terminal":
			totalArgs = append(totalArgs, "--wait", "--", "sh", "-c", strings.Join(args, " "))
		case "xterm":
			command := fmt.Sprintf("sh -c '%s'", strings.Join(args, " "))
			totalArgs = append(totalArgs, "-fs", "12", "-fa", "Unifont", "-e", command)
		default:
			command := fmt.Sprintf("sh -c '%s'", strings.Join(args, " "))
			totalArgs = append(totalArgs, "-e", command)
		}
		// 执行命令
		cmd := exec.Command(terminal, totalArgs...)
		return cmd.Run()
	}
	return fmt.Errorf("暂不支持该操作系统：%s", runtime.GOOS)
}