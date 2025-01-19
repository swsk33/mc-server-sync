package model

import (
	"encoding/json"
	"fmt"
	"gitee.com/swsk33/sclog"
)

// Result 统一的返回响应结果
type Result[T any] struct {
	// 操作是否成功
	Success bool `json:"success"`
	// 返回的消息
	Message string `json:"message,omitempty"`
	// 返回的数据，如果没有数据则为nil
	Data T `json:"data,omitempty"`
}

// CreateNullSuccessResult 创建不包含数据的成功结果
func CreateNullSuccessResult(message string) *Result[any] {
	return &Result[any]{true, message, nil}
}

// CreateSuccessResult 创建成功结果
func CreateSuccessResult[T any](message string, data T) *Result[T] {
	return &Result[T]{true, message, data}
}

// CreateFailedResult 创建失败结果
func CreateFailedResult(message string) *Result[any] {
	return &Result[any]{false, message, nil}
}

// ParseResultJson 解析Result对象的JSON结果
//
//   - resultData 结果对象的JSON数据
//
// 返回解析的内容结果
func ParseResultJson[T any](resultData []byte) (T, error) {
	// 解析JSON数据
	var resultObject Result[T]
	var zero T
	e := json.Unmarshal(resultData, &resultObject)
	if e != nil {
		sclog.ErrorLine("解析Result JSON出错！")
		return zero, e
	}
	// 查看状态
	if !resultObject.Success {
		return zero, fmt.Errorf("操作结果失败！原因：%s", resultObject.Message)
	}
	return resultObject.Data, nil
}