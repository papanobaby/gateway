package sdk

import "net/http"

func init() {
	// 像 pprof 一样引包即注册：import _ "github.com/papanobaby/gateway/sdk" 后网关挂到 http.DefaultServeMux。
	http.Handle(DefaultOptions().Path, Handler(DefaultOptions()))
}

// Register 将 gRPC 网关 Handler 注册到 mux 的 opts.Path（默认 "/grpc-gateway"）。
// 若已通过 import _ "github.com/papanobaby/gateway/sdk" 注册过 DefaultServeMux，可仅对自定义 mux 调用 Register。
func Register(mux *http.ServeMux, opts Options) {
	if opts.Path == "" {
		opts.Path = DefaultOptions().Path
	}
	mux.Handle(opts.Path, Handler(opts))
}
