package global

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// 全局HTTP客户端
var httpClient = &http.Client{
	// 设定不超时
	Timeout: 0,
	Transport: &http.Transport{
		// 从环境变量获取代理
		Proxy: http.ProxyFromEnvironment,
	},
}

// SendRequestWithHeader 发送HTTP请求，包含自定义的请求头
//
//   - path 请求路径，需要以/开头
//   - method 请求类型，例如 http.MethodGet http.MethodPost 等
//   - header 自定义请求头，可以传入nil
//   - body 请求体，无请求体传入nil
//
// 返回请求体内容，出现错误则返回错误对象
func SendRequestWithHeader(path, method string, header map[string]string, body []byte) ([]byte, error) {
	// 准备请求体
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	// 构建请求
	request, e := http.NewRequest(method, fmt.Sprintf("http://%s:%d%s", TotalConfig.Server.Host, TotalConfig.Server.Port, path), reader)
	if e != nil {
		return nil, e
	}
	// 加入自定义请求头
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}
	// 发起请求
	response, e := httpClient.Do(request)
	if e != nil {
		return nil, e
	}
	defer func() {
		_ = response.Body.Close()
	}()
	// 读取响应体
	return io.ReadAll(response.Body)
}

// SendRequest 发送HTTP请求
//
//   - path 请求路径，需要以/开头
//   - method 请求类型，例如 http.MethodGet http.MethodPost 等
//   - body 请求体，无请求体传入nil
//
// 返回请求体内容，出现错误则返回错误对象
func SendRequest(path, method string, body []byte) ([]byte, error) {
	return SendRequestWithHeader(path, method, nil, body)
}