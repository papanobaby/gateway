package sdk

import "time"

// Options 网关 SDK 的配置（可选）。
type Options struct {
	// Timeout 单次 gRPC 调用的超时，零值表示不设超时。
	Timeout time.Duration
	// Path 注册到 mux 的路径，默认 "/grpc-gateway"。
	Path string
}

// DefaultOptions 返回默认配置。
func DefaultOptions() Options {
	return Options{
		Path: "/grpc-gateway",
	}
}
