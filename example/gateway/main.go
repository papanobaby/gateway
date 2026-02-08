// 示例：像 pprof 一样引包即注册，通过 POST /grpc-gateway 调用 example 的 gRPC 服务。
// 描述符由 make gen 生成到 SDK 的 core 目录，随 SDK 分发，运行网关时从 core 包所在目录加载。
package main

import (
	"log"
	"net/http"

	_ "github.com/papanobaby/gateway/sdk"
)

func main() {
	log.Println("HTTP gateway listening on :8080 (POST /grpc-gateway)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
