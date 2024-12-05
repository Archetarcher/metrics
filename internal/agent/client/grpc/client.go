package grpc

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	pb "github.com/Archetarcher/metrics.git/proto"
)

type MetricsClient struct {
	client pb.MetricsClient

	config *config.AppConfig
}

func NewMetricsClient(config *config.AppConfig, client pb.MetricsClient) *MetricsClient {
	return &MetricsClient{
		client: client,
		config: config,
	}
}

// Run starts grpc client
func Run(c *config.AppConfig) error {

	//conn, err := grpc.Dial(c.GRPCRunAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	logger.Log.Error("failed to init grpc connection with server", zap.Error(err))
	//	return err
	//}
	//defer conn.Close()
	//
	//client := pb.NewMetricsClient(conn)
	//ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	//defer cancel()
	//
	//mc := NewMetricsClient(c, client)

	return nil
}

func (c *MetricsClient) SendMetrics(ctx context.Context) error {

	return nil
}
