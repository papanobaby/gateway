package sdk

import (
	"encoding/json"
	"net/http"

	"github.com/papanobaby/gateway/core"
)

// HTTP 请求 body 的 JSON 结构。支持 target / target_addr、method / full_method_name。
type gatewayRequest struct {
	Target         string          `json:"target"`            // gRPC 目标地址，如 "host:port"
	TargetAddr     string          `json:"target_addr"`       // 同上，兼容字段
	FullMethodName string          `json:"method"`            // 完整方法名，如 "/pb_playercenter.PlayerCenterSrv/SetWantedRole"
	FullMethodNameAlt string       `json:"full_method_name"`  // 同上，兼容字段
	Body           json.RawMessage `json:"body"`              // 请求体的 JSON
}

type errorResponse struct {
	Error string `json:"error"`
}

// Handler 返回网关的 http.Handler，描述符从 SDK 的 core 包所在目录读取（随 SDK 分发，接入方无需生成）。
func Handler(opts Options) http.Handler {
	inv := core.NewInvoker(core.DefaultDescriptorDir(), opts.Timeout)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONError(w, http.StatusMethodNotAllowed, "method must be POST")
			return
		}
		var req gatewayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
			return
		}
		if req.Target == "" {
			req.Target = req.TargetAddr
		}
		if req.FullMethodName == "" {
			req.FullMethodName = req.FullMethodNameAlt
		}
		if req.Target == "" {
			writeJSONError(w, http.StatusBadRequest, "missing target")
			return
		}
		if req.FullMethodName == "" {
			writeJSONError(w, http.StatusBadRequest, "missing method (full_method_name)")
			return
		}
		if req.Body == nil {
			req.Body = []byte("{}")
		}

		resp, err := inv.Invoke(r.Context(), &core.InvokeRequest{
			Target:         req.Target,
			FullMethodName: req.FullMethodName,
			Body:           req.Body,
		})
		if err != nil {
			writeJSONError(w, http.StatusBadGateway, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	})
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
