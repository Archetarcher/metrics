package grpc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/encryption"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

// MetricsInterceptor interceptor for grpc client.
type MetricsInterceptor struct {
	config *config.AppConfig
}

// NewMetricsInterceptor creates instance of MetricsInterceptor.
func NewMetricsInterceptor(appConfig *config.AppConfig) *MetricsInterceptor {
	return &MetricsInterceptor{config: appConfig}
}

// LoggerInterceptor logs requests.
func (i *MetricsInterceptor) LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	response, err := handler(ctx, req)
	duration := time.Since(start)

	logger.Log.Info("got incoming grpc request",
		zap.String("method", info.FullMethod),

		zap.Any("response", response),
		zap.Duration("duration", duration),
	)
	return response, err
}

// HashInterceptor hashes request.
func (i *MetricsInterceptor) HashInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var hash string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("HashSHA256g")
		if len(values) > 0 {
			hash = values[0]
		}
	}

	if i.config.Key != "" && hash != "" {
		h := hmac.New(sha256.New, []byte(i.config.Key))
		js, err := json.Marshal(req)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "incorrect hash")
		}

		h.Write(js)
		sign := h.Sum(nil)

		s, err := hex.DecodeString(hash)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "incorrect hash")
		}

		if !hmac.Equal(s, sign) {
			return nil, status.Error(codes.PermissionDenied, "incorrect hash")
		}

	}

	return handler(ctx, req)
}

// TrustedSubnetInterceptor allows only trusted ip addresses.
func (i *MetricsInterceptor) TrustedSubnetInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var ip string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("X-Real-IP")
		if len(values) > 0 {
			ip = values[0]
		}
	}

	if i.config.TrustedSubnet != "" {

		if ip == "" || ip != i.config.TrustedSubnet {
			return nil, status.Error(codes.PermissionDenied, "ip not allowed")
		}

	}
	return handler(ctx, req)
}

// DecryptInterceptor decrypts request body.
func (i *MetricsInterceptor) DecryptInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if strings.Contains(info.FullMethod, "StartSession") {
		return handler(ctx, req)
	}

	if i.config.PrivateKeyPath != "" {
		m := req.(*proto.UpdateMetricsRequest)

		decrypted, eErr := encryption.NewSymmetric(i.config.Session).Decrypt(m.Metrics)
		if eErr != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to decrypt")
		}
		req = &proto.UpdateMetricsRequest{Metrics: decrypted}

	}

	return handler(ctx, req)
}
