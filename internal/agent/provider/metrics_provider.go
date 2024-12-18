package provider

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/client/rest"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type MetricsProvider struct {
	config *config.AppConfig
	client *resty.Client
}

func NewMetricsProvider(appConfig *config.AppConfig, client *resty.Client) *MetricsProvider {
	return &MetricsProvider{config: appConfig, client: client}
}

func (p *MetricsProvider) Update(request []domain.Metrics) (*domain.SendResponse, *domain.MetricsError) {

	url := "http://" + p.config.ServerRunAddr + "/updates/"

	res, err := p.client.
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			return rest.HashMiddleware(client, request, p.config)
		}).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			return rest.GzipMiddleware(client, request, p.config)
		}).
		R().
		SetHeader("X-Real-IP", config.GetLocalIP().String()).
		SetBody(request).
		Post(url)
	if err != nil {
		return nil, &domain.MetricsError{Text: fmt.Sprintf("client: could not create request: %s\n", err.Error()), Code: http.StatusInternalServerError}
	}

	if res.StatusCode() != http.StatusOK {
		return nil, &domain.MetricsError{Text: fmt.Sprintf("client: responded with error: %s\n, %s, ", err, url), Code: res.StatusCode()}
	}
	return &domain.SendResponse{Status: http.StatusOK}, nil
}

func (p *MetricsProvider) StartSession(retryCount int) *domain.MetricsError {
	url := "http://" + p.config.ServerRunAddr + "/session/"

	key, gErr := encryption.GenKey(16)
	if gErr != nil {
		return &domain.MetricsError{Code: http.StatusInternalServerError, Text: "failed to generate crypto key"}
	}
	encryptedKey, eErr := encryption.NewAsymmetric(p.config.PublicKeyPath).Encrypt(key)
	if eErr != nil {
		return eErr
	}

	res, err := p.client.
		R().
		SetBody(domain.SessionRequest{Key: encryptedKey}).
		SetHeader("X-Real-IP", config.GetLocalIP().String()).
		Post(url)

	if (err != nil || res.StatusCode() != http.StatusOK) && retryCount > 0 {
		time.Sleep(time.Duration(p.config.ReportInterval) * time.Second)
		return p.StartSession(retryCount - 1)
	}

	if err != nil {
		return &domain.MetricsError{Text: fmt.Sprintf("client: could not create request: %s\n", err.Error()), Code: http.StatusInternalServerError}
	}

	if res.StatusCode() != http.StatusOK {
		return &domain.MetricsError{Text: fmt.Sprintf("client: responded with error creating session: %s\n, %s, %s", err, url, string(key)), Code: res.StatusCode()}
	}

	p.config.Session.Key = string(key)
	return nil
}
