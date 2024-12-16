package interceptors

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TrustedSubnetInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, config *config.AppConfig) (interface{}, error) {
	var ip string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("X-Real-IP")
		if len(values) > 0 {
			ip = values[0]
		}
	}
	if config.TrustedSubnet != "" {

		if ip == "" || ip != config.TrustedSubnet {
			return nil, status.Error(codes.PermissionDenied, "ip not allowed")
		}

	}
	return handler(ctx, req)
}
