package grpc

import (
	"bytes"
	"context"
	"encoding/json"
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
	"io"
	"net"
)

// MetricsServer is a struct for grpc server.
type MetricsServer struct {
	pb.UnimplementedMetricsServer

	service MetricsService
	config  *config.AppConfig
}

// NewMetricsServer creates instance of MetricsServer.
func NewMetricsServer(c *config.AppConfig, s MetricsService) *MetricsServer {
	return &MetricsServer{
		service: s,
		config:  c,
	}
}

// Run starts grpc server.
func (s *MetricsServer) Run() error {
	listen, err := net.Listen("tcp", s.config.GRPCRunAddr)
	if err != nil {
		logger.Log.Error("failed to define grpc port server", zap.Error(err))
		return err
	}

	interceptors := NewMetricsInterceptor(s.config)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.LoggerInterceptor,
		interceptors.HashInterceptor,
		interceptors.TrustedSubnetInterceptor,
		interceptors.DecryptInterceptor,
	))
	pb.RegisterMetricsServer(server, &MetricsServer{config: s.config, service: s.service})

	logger.Log.Info("Running grpc server ", zap.String("address", s.config.GRPCRunAddr))

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

// UpdateMetrics grpc handler for update metrics request.
func (s *MetricsServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.Empty, error) {

	var metrics []domain.Metrics

	dec := json.NewDecoder(io.NopCloser(bytes.NewReader(in.Metrics)))

	if err := dec.Decode(&metrics); err != nil {
		return nil, status.Error(codes.Internal, "cannot decode request JSON body")
	}

	_, err := s.service.Updates(ctx, metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Text, err.Code)
	}

	return &pb.Empty{}, nil

}

// StartSession grpc handler for start session request.
func (s *MetricsServer) StartSession(ctx context.Context, in *pb.StartSessionRequest) (*pb.Empty, error) {

	key, eErr := encryption.NewAsymmetric(s.config.PrivateKeyPath).Decrypt(in.Key)
	if eErr != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed")
	}

	s.config.Session = string(key)

	return &pb.Empty{}, nil

}
