package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/types"
	pb "github.com/Archetarcher/metrics.git/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type MetricsClient struct {
	client pb.MetricsClient

	service MetricsService
	config  *config.AppConfig
}

func NewMetricsClient(config *config.AppConfig, service MetricsService) *MetricsClient {
	return &MetricsClient{
		config:  config,
		service: service,
	}
}

// MetricsService is an interface for tracking metrics, sends and fetch memory and runtime metrics.
type MetricsService interface {
	TrackMetrics(ctx context.Context, update types.UpdateMetrics, group *sync.WaitGroup)
}

// Run starts grpc client
func (c *MetricsClient) Run() error {
	interceptor := newMetricsInterceptor(c.config)
	conn, err := grpc.NewClient(c.config.GRPCRunAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.hashInterceptor,
			interceptor.trustedSubnetInterceptor,
			interceptor.encryptInterceptor,
		))

	if err != nil {
		logger.Log.Error("failed to init grpc connection with server", zap.Error(err))
		return err
	}
	defer func() {
		cErr := conn.Close()
		if cErr != nil {
			log.Fatal("failed to close grpc connection")
		}
	}()

	c.client = pb.NewMetricsClient(conn)

	sErr := c.StartSession()
	if sErr != nil {
		logger.Log.Error("failed to start secure session", zap.String("error", sErr.Text), zap.Int("code", sErr.Code))
		return errors.New(sErr.Text)
	}

	c.TrackMetrics()
	return nil
}

func (c *MetricsClient) StartSession() *domain.MetricsError {
	key, gErr := encryption.GenKey(16)
	if gErr != nil {
		return &domain.MetricsError{Text: "failed to generate crypto key"}
	}

	encryptedKey, eErr := encryption.NewAsymmetric(c.config.PublicKeyPath).Encrypt(key)
	if eErr != nil {
		return eErr
	}

	_, err := c.client.StartSession(context.Background(), &pb.StartSessionRequest{Key: encryptedKey})
	if err != nil {
		return &domain.MetricsError{Text: fmt.Sprintf("client: responded with error: %s\n ", err)}
	}
	c.config.Session.Key = string(key)

	return nil
}

func (c *MetricsClient) TrackMetrics() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)

	c.service.TrackMetrics(ctx, c.update, &wg)
	logger.Log.Info("Waiting for goroutines to finish...")

	configShutdown(cancelFunc, &wg)
}

func (c *MetricsClient) update(request []domain.Metrics) (*domain.SendResponse, *domain.MetricsError) {

	js, err := json.Marshal(request)

	if err != nil {
		return nil, &domain.MetricsError{Text: fmt.Sprintf("failed to marshal request %s\n", err)}
	}

	_, cErr := c.client.UpdateMetrics(context.Background(), &pb.UpdateMetricsRequest{Metrics: js}, grpc.UseCompressor(gzip.Name))
	if cErr != nil {
		return nil, &domain.MetricsError{Text: fmt.Sprintf("client: responded with error: %s\n", cErr)}
	}
	return &domain.SendResponse{}, nil
}

func configShutdown(cancelFunc context.CancelFunc, group *sync.WaitGroup) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigint

	logger.Log.Info("Shutdown signal received")

	cancelFunc()
	group.Wait()
}
