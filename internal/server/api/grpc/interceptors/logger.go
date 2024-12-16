package interceptors

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	response, err := handler(ctx, req)
	duration := time.Since(start)

	logger.Log.Info("got incoming HTTP request",
		zap.String("method", info.FullMethod),

		zap.Any("response", response),
		zap.Duration("duration", duration),
	)
	return response, err
}
