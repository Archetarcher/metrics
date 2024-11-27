package handlers

import (
	"context"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/server/encryption"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
)

const emptyParam = ""

// MetricsHandler is a handler to work with metrics, keeps implementation of MetricsService and configuration
type MetricsHandler struct {
	service MetricsService
	config  *config.AppConfig
}

// MetricsService is an interface that describes interaction with service layer
type MetricsService interface {
	Updates(request []domain.Metrics, ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	Update(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError)
	GetValue(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError)
	GetAllValues(ctx context.Context) (string, *domain.MetricsError)
	CheckConnection(ctx context.Context) *domain.MetricsError
}

// NewMetricsHandler creates new handler
func NewMetricsHandler(service MetricsService, appConfig *config.AppConfig) *MetricsHandler {
	return &MetricsHandler{service: service, config: appConfig}
}

// StartSession handler that accepts and saves session key.
func (h *MetricsHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateSessionRequest(r)

	enc := json.NewEncoder(w)

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	key, eErr := encryption.DecryptAsymmetric(request.Key, h.config.PrivateKeyPath)
	if eErr != nil {
		sendResponse(enc, eErr.Error(), http.StatusUnauthorized, w)
		return
	}

	h.config.Session = string(key)
	sendResponse(enc, "", http.StatusOK, w)
}

// UpdateMetrics handler that creates or updates existing metric.
// Data provided in path parameters
func (h *MetricsHandler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateRequest(r)

	enc := json.NewEncoder(w)

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	_, err = h.service.Update(request, r.Context())
	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}
	sendResponse(enc, "", http.StatusOK, w)
}

// GetMetrics handler that returns existing metric by ID and MType in domain.Metrics.
// Data provided in path parameters.
func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateGetRequest(r)

	enc := json.NewEncoder(w)

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	result, err := h.service.GetValue(request, r.Context())
	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	var resp any

	if result.MType == domain.CounterType {
		resp = result.Delta
	}
	if result.MType == domain.GaugeType {
		resp = result.Value
	}
	sendResponse(enc, resp, http.StatusOK, w)
}

// UpdatesMetrics handler that creates or updates existing batch of metrics.
// Data provided in body json.
func (h *MetricsHandler) UpdatesMetrics(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	// validate
	request, err := validateUpdatesRequest(r)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	_, err = h.service.Updates(request, r.Context())
	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	sendResponse(enc, "", http.StatusOK, w)
}

// UpdateMetricsJSON handler that creates or updates existing metric.
// Data provided in body json format.
func (h *MetricsHandler) UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateRequest(r)

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	response, err := h.service.Update(request, r.Context())
	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	sendResponse(enc, response, http.StatusOK, w)
}

// GetMetricsJSON handler that returns existing metric by ID and MType in domain.Metrics.
// Data provided in body json format
func (h *MetricsHandler) GetMetricsJSON(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateGetRequest(r)

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	response, err := h.service.GetValue(request, r.Context())
	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	sendResponse(enc, response, http.StatusOK, w)

}

// GetMetricsPage handler that returns all metrics from database, in table view.
func (h *MetricsHandler) GetMetricsPage(w http.ResponseWriter, r *http.Request) {

	result, err := h.service.GetAllValues(r.Context())
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	sendResponse(enc, result, http.StatusOK, w)
}

// GetPing handler that checks database connection.
func (h *MetricsHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	err := h.service.CheckConnection(r.Context())

	if err != nil {
		sendResponse(enc, err.Text, http.StatusInternalServerError, w)
		return
	}

	sendResponse(enc, "", http.StatusOK, w)

}

func validateGetRequest(r *http.Request) (*domain.Metrics, *domain.MetricsError) {

	// validate params
	var metrics domain.Metrics

	n := chi.URLParam(r, "name")
	t := chi.URLParam(r, "type")

	metrics.ID = n
	metrics.MType = t

	if metrics.ID == emptyParam && metrics.MType == emptyParam {
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&metrics); err != nil {
			return nil, &domain.MetricsError{
				Text: "cannot decode request JSON body",
				Code: http.StatusInternalServerError,
			}
		}
	}

	if metrics.ID == emptyParam || metrics.MType == emptyParam {
		return nil, &domain.MetricsError{
			Text: "empty param",
			Code: http.StatusBadRequest,
		}
	}

	if !slices.Contains([]string{domain.GaugeType, domain.CounterType}, metrics.MType) {
		return nil, &domain.MetricsError{
			Text: "incorrect type",
			Code: http.StatusBadRequest,
		}
	}

	return &metrics, nil
}
func validateSessionRequest(r *http.Request) (*domain.SessionRequest, *domain.MetricsError) {
	// validate method
	if r.Method != http.MethodPost {

		return nil, &domain.MetricsError{
			Text: "method not allowed",
			Code: http.StatusMethodNotAllowed,
		}
	}

	// validate params
	var session domain.SessionRequest

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&session); err != nil {
		return nil, &domain.MetricsError{
			Text: "cannot decode request JSON body",
			Code: http.StatusInternalServerError,
		}
	}

	return &session, nil
}
func validateRequest(r *http.Request) (*domain.Metrics, *domain.MetricsError) {
	// validate method
	if r.Method != http.MethodPost {

		return nil, &domain.MetricsError{
			Text: "method not allowed",
			Code: http.StatusMethodNotAllowed,
		}
	}

	// validate params
	var metrics domain.Metrics

	n := chi.URLParam(r, "name")
	t := chi.URLParam(r, "type")
	v := chi.URLParam(r, "value")

	metrics.ID = n
	metrics.MType = t

	if t == domain.GaugeType {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, &domain.MetricsError{
				Text: "",
				Code: http.StatusBadRequest,
			}
		}
		metrics.Value = &value
	}

	if t == domain.CounterType {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, &domain.MetricsError{
				Text: "",
				Code: http.StatusBadRequest,
			}
		}
		metrics.Delta = &value
	}

	if metrics.ID == emptyParam && metrics.MType == emptyParam && metrics.Delta == nil && metrics.Value == nil {
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&metrics); err != nil {
			return nil, &domain.MetricsError{
				Text: "cannot decode request JSON body",
				Code: http.StatusInternalServerError,
			}
		}
	}

	if metrics.ID == emptyParam || metrics.MType == emptyParam || ((metrics.MType == domain.GaugeType && metrics.Value == nil) || (metrics.MType == domain.CounterType && metrics.Delta == nil)) {
		return nil, &domain.MetricsError{
			Text: "empty param",
			Code: http.StatusBadRequest,
		}
	}

	if !slices.Contains([]string{domain.GaugeType, domain.CounterType}, metrics.MType) {
		return nil, &domain.MetricsError{
			Text: "incorrect type",
			Code: http.StatusBadRequest,
		}
	}

	return &metrics, nil
}
func validateUpdatesRequest(r *http.Request) ([]domain.Metrics, *domain.MetricsError) {
	// validate method
	if r.Method != http.MethodPost {

		return nil, &domain.MetricsError{
			Text: "method not allowed",
			Code: http.StatusMethodNotAllowed,
		}
	}

	// validate params
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, &domain.MetricsError{
			Text: err.Error(),
			Code: http.StatusBadRequest,
		}
	}

	metrics := make([]domain.Metrics, 3)
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		return nil, &domain.MetricsError{
			Text: err.Error(),
			Code: http.StatusInternalServerError,
		}
	}

	for _, m := range metrics {
		if m.MType == emptyParam ||
			!slices.Contains([]string{domain.GaugeType, domain.CounterType}, m.MType) ||
			m.ID == emptyParam {
			return nil, &domain.MetricsError{
				Text: "invalid param provided",
				Code: http.StatusBadRequest,
			}
		}
	}

	return metrics, nil
}

func sendResponse(enc *json.Encoder, data interface{}, code int, w http.ResponseWriter) {
	w.WriteHeader(code)

	if code > http.StatusOK {
		logger.Log.Info("failed with error", zap.Any("error", data), zap.Int("code", code))
	}
	if err := enc.Encode(data); err != nil {
		logger.Log.Debug("error encryption response", zap.Error(err))
	}

}
