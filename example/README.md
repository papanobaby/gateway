# 示例：gRPC 服务 + HTTP 网关

描述符生成到 **SDK 的 core 目录**（example 的 make gen 会写到 ../core/），随 SDK 分发，接入方无需各自生成；网关引包即注册。

## 生成描述符与 Go 代码

需安装 `protoc`、`protoc-gen-go`、`protoc-gen-go-grpc`，并将 `$(go env GOPATH)/bin` 加入 PATH：

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# 在 example 目录下，描述符生成到 ../core/，Go 代码到 pb/
PATH="$(go env GOPATH)/bin:$PATH" make gen
```

## 运行

1. 启动 gRPC 服务（终端一）：
   ```bash
   cd example && go run ./grpc_server
   ```

2. 启动 HTTP 网关（终端二）：
   ```bash
   cd example && go run ./gateway
   ```
   （描述符从 SDK 的 core 包目录加载，与运行目录无关。）

3. 通过网关调用 gRPC：
   ```bash
   curl -X POST http://localhost:8080/grpc-gateway \
     -H "Content-Type: application/json" \
     -d '{"target":"127.0.0.1:50051","method":"/echo.EchoService/Echo","body":{"message":"hello"}}'
   ```
   预期返回：`{"message":"echo: hello"}`
