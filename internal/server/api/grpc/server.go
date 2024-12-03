package grpc

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/encryption"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	pb "github.com/Archetarcher/metrics.git/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"slices"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	repo   MetricRepository
	config *config.AppConfig
}

func RunGRPCServer(c *config.AppConfig, r MetricRepository) error {
	listen, err := net.Listen("tcp", c.GRPCRunAddr)
	if err != nil {
		logger.Log.Error("failed to define grpc port server", zap.Error(err))
		return err
	}

	s := grpc.NewServer()
	pb.RegisterMetricsServer(s, &MetricsServer{config: c, repo: r})

	logger.Log.Info("Running grpc server ", zap.String("address", c.GRPCRunAddr))

	if sErr := s.Serve(listen); sErr != nil {
		logger.Log.Error("failed to serve grpc server", zap.Error(sErr))
		return sErr
	}
	return nil
}

// MetricRepository is an interface that describes interaction with repository layer
type MetricRepository interface {
	GetAllIn(ctx context.Context, keys []string) ([]domain.Metrics, *domain.MetricsError)
	SetAll(ctx context.Context, request []domain.Metrics) *domain.MetricsError
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

	keys := make([]string, len(metrics))
	for _, m := range metrics {
		if !slices.Contains(keys, getKey(m)) {
			keys = append(keys, getKey(m))
		}
	}

	metricsByKey, err := s.repo.GetAllIn(ctx, keys)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Text, err.Code)
	}

	for _, mbk := range metricsByKey {
		for k, m := range metrics {
			if getKey(m) == getKey(mbk) && m.MType == domain.CounterType {
				c := *m.Delta + *mbk.Delta
				metrics[k].Delta = &c
			}
		}
	}

	if rErr := s.repo.SetAll(ctx, metrics); rErr != nil {
		return nil, status.Errorf(codes.Internal, rErr.Text, rErr.Code)
	}

	return &pb.Empty{}, nil

}

func (s *MetricsServer) StartSession(ctx context.Context, in *pb.StartSessionRequest) (*pb.Empty, error) {

	key, eErr := encryption.DecryptAsymmetric(in.Key, s.config.PrivateKeyPath)
	if eErr != nil {
		return nil, status.Errorf(codes.Unauthenticated, eErr.Error())
	}

	s.config.Session = string(key)
	return &pb.Empty{}, nil

}

func getKey(request domain.Metrics) string {
	return request.ID + "_" + request.MType
}
