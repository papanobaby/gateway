package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// DefaultDescriptorDir 返回 core 包所在目录（描述符 .pb 文件存放于此，随 SDK 分发，接入方无需各自生成）。
func DefaultDescriptorDir() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Dir(f)
}

// ParseFullMethodName 从 gRPC 完整方法名 "/package.Service/Method" 解析出服务名 "package.Service"。
func ParseFullMethodName(fullMethodName string) (serviceName string, methodName string, err error) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/")
	parts := strings.SplitN(fullMethodName, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid full method name: %q, expected /package.Service/Method", fullMethodName)
	}
	return parts[0], parts[1], nil
}

// MethodResolver 根据 full_method_name 解析并缓存 *desc.MethodDescriptor。
type MethodResolver struct {
	descriptorDir string
	mu           sync.RWMutex
	cache        map[string]*desc.MethodDescriptor
}

// NewMethodResolver 创建方法描述符解析器，descriptorDir 为存放 .pb 文件的目录。
func NewMethodResolver(descriptorDir string) *MethodResolver {
	return &MethodResolver{
		descriptorDir: descriptorDir,
		cache:        make(map[string]*desc.MethodDescriptor),
	}
}

// Resolve 根据 gRPC 完整方法名（如 "/pb_playercenter.PlayerCenterSrv/SetWantedRole"）返回 *desc.MethodDescriptor。
func (r *MethodResolver) Resolve(fullMethodName string) (*desc.MethodDescriptor, error) {
	r.mu.RLock()
	md, ok := r.cache[fullMethodName]
	r.mu.RUnlock()
	if ok {
		return md, nil
	}

	serviceName, _, err := ParseFullMethodName(fullMethodName)
	if err != nil {
		return nil, err
	}

	// 约定：描述符文件名为 {服务名}.pb，与 full_method_name 中的服务名一致
	pbPath := filepath.Join(r.descriptorDir, serviceName+".pb")
	data, err := os.ReadFile(pbPath)
	if err != nil {
		return nil, fmt.Errorf("read descriptor file %s: %w", pbPath, err)
	}

	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &fds); err != nil {
		return nil, fmt.Errorf("unmarshal FileDescriptorSet: %w", err)
	}

	files, err := desc.CreateFileDescriptorsFromSet(&fds)
	if err != nil {
		return nil, fmt.Errorf("create file descriptors: %w", err)
	}

	// 构建 gRPC 完整方法名格式：/package.Service/Method
	for _, fd := range files {
		for _, svc := range fd.GetServices() {
			for _, m := range svc.GetMethods() {
				fqn := "/" + svc.GetFullyQualifiedName() + "/" + m.GetName()
				if fqn == fullMethodName {
					r.mu.Lock()
					r.cache[fullMethodName] = m
					r.mu.Unlock()
					return m, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("method %q not found in descriptor set", fullMethodName)
}
