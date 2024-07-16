package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"slices"
	"strconv"
	"time"
)

type MetricsHandler struct {
	service MetricsService
	config  *config.AppConfig
}

type MetricsService interface {
	Update(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError)
	GetValue(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError)
	GetAllValues() (string, *domain.MetricsError)
}

func NewMetricsHandler(service MetricsService, appConfig *config.AppConfig) *MetricsHandler {
	return &MetricsHandler{service: service, config: appConfig}
}

func (h *MetricsHandler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateRequest(r)

	enc := json.NewEncoder(w)

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	_, err = h.service.Update(request)
	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}
	sendResponse(enc, "", http.StatusOK, w)
}
func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateGetRequest(r)

	enc := json.NewEncoder(w)

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	result, err := h.service.GetValue(request)
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

func (h *MetricsHandler) UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateRequest(r)

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	response, err := h.service.Update(request)
	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	sendResponse(enc, response, http.StatusOK, w)
}
func (h *MetricsHandler) GetMetricsJSON(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateGetRequest(r)

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	response, err := h.service.GetValue(request)
	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	sendResponse(enc, response, http.StatusOK, w)

}
func (h *MetricsHandler) GetMetricsPage(w http.ResponseWriter, r *http.Request) {

	result, err := h.service.GetAllValues()
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err != nil {
		sendResponse(enc, err.Text, err.Code, w)
		return
	}

	sendResponse(enc, result, http.StatusOK, w)
}

func (h *MetricsHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("pgx", h.config.Store.Pgx.DatabaseDsn)

	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = db.PingContext(ctx)

	if err != nil {
		sendResponse(enc, err.Error(), http.StatusInternalServerError, w)
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

	if metrics.ID == domain.EmptyParam && metrics.MType == domain.EmptyParam {
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&metrics); err != nil {
			return nil, &domain.MetricsError{
				Text: "cannot decode request JSON body",
				Code: http.StatusInternalServerError,
			}
		}
	}

	if metrics.ID == domain.EmptyParam || metrics.MType == domain.EmptyParam {
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

	if metrics.ID == domain.EmptyParam && metrics.MType == domain.EmptyParam && metrics.Delta == nil && metrics.Value == nil {
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&metrics); err != nil {
			return nil, &domain.MetricsError{
				Text: "cannot decode request JSON body",
				Code: http.StatusInternalServerError,
			}
		}
	}

	if metrics.ID == domain.EmptyParam || metrics.MType == domain.EmptyParam || ((metrics.MType == domain.GaugeType && metrics.Value == nil) || (metrics.MType == domain.CounterType && metrics.Delta == nil)) {
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

func sendResponse(enc *json.Encoder, data interface{}, code int, w http.ResponseWriter) {
	w.WriteHeader(code)

	if code > http.StatusOK {
		logger.Log.Info("failed with error", zap.Any("error", data), zap.Int("code", code))
	}
	if err := enc.Encode(data); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
	}

}
