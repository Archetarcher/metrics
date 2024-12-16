package interceptors

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func HashInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, config *config.AppConfig) (interface{}, error) {
	var hash string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("HashSHA256")
		if len(values) > 0 {
			hash = values[0]
		}
	}

	if config.Key != "" && hash != "" {
		h := hmac.New(sha256.New, []byte(config.Key))
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
