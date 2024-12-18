package grpc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"github.com/Archetarcher/metrics.git/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

type metricsInterceptor struct {
	config *config.AppConfig
}

func newMetricsInterceptor(appConfig *config.AppConfig) *metricsInterceptor {
	return &metricsInterceptor{config: appConfig}
}

func (i *metricsInterceptor) hashInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	if i.config.Key != "" {
		h := hmac.New(sha256.New, []byte(i.config.Key))

		js, err := json.Marshal(req)

		if err != nil {
			return err
		}

		h.Write(js)
		hash := h.Sum(nil)

		md := metadata.Pairs("HashSHA256g", hex.EncodeToString(hash))
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func (i *metricsInterceptor) trustedSubnetInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	md := metadata.New(map[string]string{"X-Real-IP": config.GetLocalIP().String()})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	return invoker(ctx, method, req, reply, cc, opts...)
}

func (i *metricsInterceptor) encryptInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	if strings.Contains(method, "StartSession") {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	m := req.(*proto.UpdateMetricsRequest)

	encrypted := encryption.NewSymmetric(i.config.Session.Key).Encrypt(m.Metrics)

	req = &proto.UpdateMetricsRequest{Metrics: encrypted}

	return invoker(ctx, method, req, reply, cc, opts...)
}
