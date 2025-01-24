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

// ModInfoMap 模组列表哈希表类型，其中：
//
//   - 键：模组文件的sha256摘要
//   - 值：完整模组文件元数据
type ModInfoMap map[string]*ModFile

// NewModInfoMapFromSlice 从模组列表新建模组哈希表
//
//   - list 模组列表
//
// 返回新的模组哈希表对象
func NewModInfoMapFromSlice(list []*ModFile) ModInfoMap {
	modMap := make(ModInfoMap)
	for _, item := range list {
		modMap[item.Sha256] = item
	}
	return modMap
}

// FindDuplicate 查找模组列表切片中SHA256重复的模组对象
//
//   - list 要查找的列表
//
// 返回重复的模组对象
func FindDuplicate(list []*ModFile) []*ModFile {
	duplicates := make([]*ModFile, 0)
	modMap := make(ModInfoMap)
	for _, item := range list {
		if modMap.Contains(item) {
			duplicates = append(duplicates, item)
			continue
		}
		modMap[item.Sha256] = item
	}
	return duplicates
}

// Contains 查看当前模组哈希表中是否存在某个模组
//
//   - mod 要判断存在的模组对象
//
// 如果mod存在于当前modMap中，返回true，否则返回false
func (modMap ModInfoMap) Contains(mod *ModFile) bool {
	_, ok := modMap[mod.Sha256]
	return ok
}

// Remove 从modMap移除一些元素
//
//   - removeMods 要移除的模组对象，为不定长参数
func (modMap ModInfoMap) Remove(removeMods ...*ModFile) {
	for _, item := range removeMods {
		delete(modMap, item.Sha256)
	}
}

// Subtract 哈希表作差，查找出当前哈希表modMap存在但是传入哈希表subtractMap不存在的模组信息对象
//
//   - subtractMap 传入的作差哈希表
//
// 返回当前哈希表modMap存在但是传入哈希表subtractMap不存在的 ModFile 对象指针切片
func (modMap ModInfoMap) Subtract(subtractMap ModInfoMap) []*ModFile {
	list := make([]*ModFile, 0)
	for _, item := range modMap {
		if !subtractMap.Contains(item) {
			list = append(list, item)
		}
	}
	return list
}

// ToSlice 将模组哈希表转换成模组信息对象切片
//
// 返回切片形式的模组列表对象
func (modMap ModInfoMap) ToSlice() []*ModFile {
	list := make([]*ModFile, 0)
	for _, item := range modMap {
		list = append(list, item)
	}
	return list
}