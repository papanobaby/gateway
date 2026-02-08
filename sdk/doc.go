// Package sdk 提供 gRPC 网关 HTTP 接口，类似 net/http/pprof：引包即注册。
//
// 功能：对外暴露一个 HTTP 接口，请求体中指定 gRPC 目标地址、完整方法名（如 /package.Service/Method）和 JSON body；
// SDK 从 core 包所在目录（即 SDK 代码目录）按服务名加载 .pb（FileDescriptorSet），完成 JSON 与 PB 的转换后向目标发起 gRPC 调用，并将响应转为 JSON 返回。
//
// 描述符：存放在 SDK 的 core 包目录下，随 SDK 一起分发，接入方无需各自生成。文件名 = {服务全限定名}.pb，
// 例如 /pb_playercenter.PlayerCenterSrv/SetWantedRole 对应 core/pb_playercenter.PlayerCenterSrv.pb。
//
// 接入示例（引包即可，像 pprof）：
//
//	import _ "github.com/papanobaby/gateway/sdk"
//
//	func main() {
//	    http.ListenAndServe(":8080", nil) // 网关已挂在 http.DefaultServeMux，路径 /grpc-gateway
//	}
//
// 自定义 mux 时：sdk.Register(mux, sdk.Options{Timeout: 10*time.Second})
//
// 请求格式：POST，Content-Type: application/json，body 示例：
//
//	{"target": "host:port", "method": "/pb_playercenter.PlayerCenterSrv/SetWantedRole", "body": {"key": "value"}}
//
// 安全：该接口应仅在内网或受控环境使用，不要暴露到公网。
package sdk
