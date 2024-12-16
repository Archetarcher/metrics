package grpc

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/api/grpc/interceptors"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/encryption"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	pb "github.com/Archetarcher/metrics.git/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
	"net"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	service MetricsService
	config  *config.AppConfig
}

func Run(c *config.AppConfig, s MetricsService) error {
	listen, err := net.Listen("tcp", c.GRPCRunAddr)
	if err != nil {
		logger.Log.Error("failed to define grpc port server", zap.Error(err))
		return err
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.LoggerInterceptor,
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			return interceptors.TrustedSubnetInterceptor(ctx, req, info, handler, c)
		},
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			return interceptors.HashInterceptor(ctx, req, info, handler, c)
		},
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			return interceptors.DecryptInterceptor(ctx, req, info, handler, c)
		}),
	)
	pb.RegisterMetricsServer(server, &MetricsServer{config: c, service: s})

	logger.Log.Info("Running grpc server ", zap.String("address", c.GRPCRunAddr))

	if sErr := server.Serve(listen); sErr != nil {
		logger.Log.Error("failed to serve grpc server", zap.Error(sErr))
		return sErr
	}
	return nil
}

// MetricsService is an interface that describes interaction with service layer
type MetricsService interface {
	Updates(ctx context.Context, request []domain.Metrics) ([]domain.Metrics, *domain.MetricsError)
}

func (s *MetricsServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.Empty, error) {

	var metrics []domain.Metrics

	for _, m := range in.Metrics {
		metrics = append(metrics, domain.Metrics{
			Delta: &m.Delta,
			Value: &m.Value,
			ID:    m.ID,
			MType: m.MType,
		})
	}
	_, err := s.service.Updates(ctx, metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Text, err.Code)
	}

	return &pb.Empty{}, nil

}
func (s *MetricsServer) StartSession(ctx context.Context, in *pb.StartSessionRequest) (*pb.Empty, error) {

	key, eErr := encryption.DecryptAsymmetric(in.Key, s.config.PrivateKeyPath)
	if eErr != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed")
	}
	logger.Log.Info("got session request:", zap.String("key", string(key)))

	s.config.Session = string(key)

	return &pb.Empty{}, nil

}
