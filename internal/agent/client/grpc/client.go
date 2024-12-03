package grpc

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	pb "github.com/Archetarcher/metrics.git/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"time"
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

	conn, err := grpc.Dial(c.GRPCRunAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("failed to init grpc connection with server", zap.Error(err))
		return err
	}
	defer conn.Close()

	client := pb.NewMetricsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	mc := NewMetricsClient(c, client)

	sErr := mc.StartSession(ctx, c.Session.RetryConn)
	if sErr != nil {
		logger.Log.Error("failed to start secure session in grpc connection with server", zap.Error(err))
		return sErr
	}

	return nil
}

func (c *MetricsClient) StartSession(ctx context.Context, retryCount int) error {
	key, gErr := encryption.GenKey(16)
	if gErr != nil {
		return gErr
	}
	encryptedKey, eErr := encryption.EncryptAsymmetric(key, c.config.PublicKeyPath)
	if eErr != nil {
		return eErr.Err
	}

	_, sErr := c.client.StartSession(ctx, &pb.StartSessionRequest{Key: encryptedKey})

	if sErr != nil {
		if e, ok := status.FromError(sErr); ok {
			if e.Code() == codes.Unauthenticated {
				logger.Log.Error("authentication failed", zap.Error(sErr))

				return sErr
			}
		}
	}

	if sErr != nil && retryCount > 0 {
		time.Sleep(time.Duration(c.config.ReportInterval) * time.Second)
		return c.StartSession(ctx, retryCount-1)
	}

	if sErr != nil {
		return sErr
	}

	c.config.Session.Key = string(key)
	return nil
}
func (c *MetricsClient) SendMetrics(ctx context.Context) error {

	return nil
}
