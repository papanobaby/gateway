package core

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

var jsonpbMarshaler = &jsonpb.Marshaler{EmitDefaults: true}

// JSONToMessage 将 JSON 请求体转为该方法输入类型的 dynamic.Message（兼容 grpcdynamic 的 proto.Message）。
func JSONToMessage(method *desc.MethodDescriptor, jsonBody []byte) (proto.Message, error) {
	msgDesc := method.GetInputType()
	msg := dynamic.NewMessage(msgDesc)
	if err := jsonpb.Unmarshal(bytes.NewReader(jsonBody), msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// MessageToJSON 将 gRPC 响应的 proto 消息转为 JSON。
func MessageToJSON(msg proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	if err := jsonpbMarshaler.Marshal(&buf, msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
