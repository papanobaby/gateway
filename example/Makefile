# 生成 Go 代码与描述符文件，需安装 protoc、protoc-gen-go、protoc-gen-go-grpc
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 描述符生成到 SDK 的 core 目录，随 SDK 分发，接入方无需再生成
PROTO_DIR := proto
CORE_DIR := ../core
PB_GO_DIR := pb

.PHONY: gen clean

gen:
	mkdir -p $(PB_GO_DIR)
	protoc -I$(PROTO_DIR) \
		--descriptor_set_out=$(CORE_DIR)/echo.EchoService.pb --include_imports \
		$(PROTO_DIR)/echo.proto
	protoc -I$(PROTO_DIR) \
		--go_out=$(PB_GO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_GO_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/echo.proto

clean:
	rm -rf $(CORE_DIR)/echo.*.pb $(PB_GO_DIR)/*.pb.go
