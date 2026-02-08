package core

import (
	"context"
	"fmt"
	"time"

	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Invoker 根据描述符目录与目标地址发起 gRPC 调用。
type Invoker struct {
	resolver *MethodResolver
	timeout  time.Duration
}

// NewInvoker 创建调用器，descriptorDir 为 .pb 所在目录，timeout 为单次 gRPC 调用超时。
func NewInvoker(descriptorDir string, timeout time.Duration) *Invoker {
	return &Invoker{
		resolver: NewMethodResolver(descriptorDir),
		timeout:  timeout,
	}
}

// InvokeRequest 为 HTTP 网关的入参。
type InvokeRequest struct {
	Target        string          // gRPC 目标地址，如 "host:port"
	FullMethodName string         // 完整方法名，如 "/pb_playercenter.PlayerCenterSrv/SetWantedRole"
	Body          []byte          // 请求体的 JSON
}

// Invoke 执行一次 Unary gRPC 调用，将 Body（JSON）转为 PB 请求，调用目标服务，响应转为 JSON 返回。
func (inv *Invoker) Invoke(ctx context.Context, req *InvokeRequest) ([]byte, error) {
	if inv.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, inv.timeout)
		defer cancel()
	}

	method, err := inv.resolver.Resolve(req.FullMethodName)
	if err != nil {
		return nil, fmt.Errorf("resolve method: %w", err)
	}

	if method.IsClientStreaming() || method.IsServerStreaming() {
		return nil, fmt.Errorf("streaming method not supported: %s", req.FullMethodName)
	}

	reqMsg, err := JSONToMessage(method, req.Body)
	if err != nil {
		return nil, fmt.Errorf("json to message: %w", err)
	}

	conn, err := grpc.DialContext(ctx, req.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", req.Target, err)
	}
	defer conn.Close()

	stub := grpcdynamic.NewStub(conn)
	respMsg, err := stub.InvokeRpc(ctx, method, reqMsg)
	if err != nil {
		return nil, fmt.Errorf("invoke rpc: %w", err)
	}

	return MessageToJSON(respMsg)
}
