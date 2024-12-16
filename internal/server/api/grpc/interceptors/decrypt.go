package interceptors

import (
	"context"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/encryption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DecryptInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, config *config.AppConfig) (interface{}, error) {
	var isEncrypted string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("Encrypted")
		if len(values) > 0 {
			isEncrypted = values[0]
		}
	}

	if config.PrivateKeyPath != "" && isEncrypted != "" {
		js, err := json.Marshal(req)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to marshal")
		}

		decrypted := encryption.DecryptSymmetric(js, config.Session)
		req = decrypted
	}

	return handler(ctx, req)
}
