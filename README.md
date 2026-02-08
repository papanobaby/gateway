# gRPC 网关 SDK

类似 `net/http/pprof`：**引包即注册**，为现有 HTTP 服务注入 gRPC 网关接口。通过 HTTP 指定目标地址、完整方法名和 JSON body，SDK 从 **core 包所在目录**（即 SDK 代码目录）加载对应服务的 .pb（FileDescriptorSet），完成 JSON 与 PB 的转换后发起 gRPC 调用并返回 JSON。

## 功能

- **HTTP 接口**：POST 请求，body 中指定 `target`（gRPC 地址）、`method`（完整方法名）、`body`（JSON 请求体）。
- **方法名格式**：gRPC 完整方法名 `/package.Service/Method`，例如 `/pb_playercenter.PlayerCenterSrv/SetWantedRole`。
- **描述符**：存放在 SDK 的 **core 包目录**下，随 SDK 一起分发，接入方无需各自生成；按服务名加载 `.pb` 文件。

## 描述符目录约定

- 描述符放在 **gateway_sdk 的 core 包目录**下（即 SDK 代码所在目录），随 SDK 分发，接入方不用每边都生成一次。
- 文件名 = **服务全限定名**.pb，与完整方法名中的服务名一致。  
  例如方法 `/pb_playercenter.PlayerCenterSrv/SetWantedRole` 对应 `core/pb_playercenter.PlayerCenterSrv.pb`。
- 在 SDK 仓库中生成并提交：`protoc --descriptor_set_out=./core/pb_xxx.MyService.pb --include_imports -I. your.proto`

## 接入方式（像 pprof，引包就行）

```go
package main

import (
	"net/http"

	_ "gateway_sdk/sdk" // 引包即注册，网关挂在 http.DefaultServeMux，路径 /grpc-gateway
)

func main() {
	http.ListenAndServe(":8080", nil)
}
```

自定义 mux 或路径时：`sdk.Register(mux, sdk.Options{Timeout: 10*time.Second, Path: "/grpc-gateway"})`

## HTTP 请求/响应

- **请求**：`POST`，`Content-Type: application/json`，body 示例：
  ```json
  {
    "target": "127.0.0.1:50051",
    "method": "/pb_playercenter.PlayerCenterSrv/SetWantedRole",
    "body": { "role_id": "123" }
  }
  ```
  兼容字段：`target_addr`、`full_method_name`。
- **成功**：`200`，body 为 gRPC 响应转成的 JSON。
- **失败**：`4xx/5xx`，body 为 `{"error": "..."}`。

## 安全说明

该接口应仅在内网或受控环境使用，不要暴露到公网；建议挂到带鉴权/白名单的路由下。

## 示例与测试

见 [example](example/) 目录：最小 gRPC 服务 + 描述符生成到 **core/** + 引包即注册的网关。  
先在该目录执行 `make gen`（需 protoc 与 Go 插件），再分别启动 `grpc_server` 与 `gateway`，通过 `POST /grpc-gateway` 调用 gRPC。
# gateway
