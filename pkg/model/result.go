package model

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