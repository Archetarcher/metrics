package handlers

import (
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"slices"
	"strconv"
)

type MetricsHandler struct {
	MetricsService
}

type MetricsService interface {
	Update(request *models.Metrics) (*models.Metrics, *models.MetricError)
	GetValue(request *models.Metrics) (*models.Metrics, *models.MetricError)
	GetAllValues() (string, *models.MetricError)
}

func (h *MetricsHandler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateRequest(r)
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Text))
		return
	}

	_, err = h.Update(request)
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Text))
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {

	// validate
	request, err := validateGetRequest(r)
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Text))
		return
	}

	result, err := h.GetValue(request)
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Text))
		return
	}

	enc := json.NewEncoder(w)

	if result.MType == models.GaugeType {
		sendResponse(enc, result.Value)
		return
	}

	sendResponse(enc, result.Delta)
	w.WriteHeader(http.StatusOK)

}

func (h *MetricsHandler) UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateRequest(r)

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(err.Code)
		sendResponse(enc, err)
		return
	}

	response, err := h.Update(request)
	if err != nil {
		w.WriteHeader(err.Code)
		sendResponse(enc, err)
		return
	}

	sendResponse(enc, response)
	w.WriteHeader(http.StatusOK)
}
func (h *MetricsHandler) GetMetricsJSON(w http.ResponseWriter, r *http.Request) {
	// validate
	request, err := validateGetRequest(r)

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(err.Code)
		sendResponse(enc, err)
		return
	}

	response, err := h.GetValue(request)
	if err != nil {
		w.WriteHeader(err.Code)
		sendResponse(enc, err)
		return
	}

	sendResponse(enc, response)
	w.WriteHeader(http.StatusOK)

}
func (h *MetricsHandler) GetMetricsPage(w http.ResponseWriter, r *http.Request) {

	result, err := h.GetAllValues()
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Text))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))

}

func validateGetRequest(r *http.Request) (*models.Metrics, *models.MetricError) {

	// validate headers
	for k, v := range models.AllowedHeaders {
		if h := r.Header.Get(k); h != v {
			return nil, &models.MetricError{
				Text: "header not allowed",
				Code: http.StatusBadRequest,
			}
		}
	}

	// validate params
	var metrics models.Metrics

	n := chi.URLParam(r, "name")
	t := chi.URLParam(r, "type")

	metrics.ID = n
	metrics.MType = t

	if reflect.DeepEqual(metrics, models.Metrics{}) {
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&metrics); err != nil {
			logger.Log.Info("cannot decode request JSON body", zap.Error(err))
			return nil, &models.MetricError{
				Text: "cannot decode request JSON body",
				Code: http.StatusInternalServerError,
			}
		}
	}

	if metrics.ID == models.EmptyParam || metrics.MType == models.EmptyParam {
		return nil, &models.MetricError{
			Text: "empty param",
			Code: http.StatusBadRequest,
		}
	}

	if !slices.Contains([]string{models.GaugeType, models.CounterType}, metrics.MType) {
		return nil, &models.MetricError{
			Text: "incorrect type",
			Code: http.StatusBadRequest,
		}
	}

	return &metrics, nil
}
func validateRequest(r *http.Request) (*models.Metrics, *models.MetricError) {
	// validate method
	if r.Method != http.MethodPost {

		return nil, &models.MetricError{
			Text: "method not allowed",
			Code: http.StatusMethodNotAllowed,
		}
	}

	// validate headers
	for k, v := range models.AllowedHeaders {
		if h := r.Header.Get(k); h != v {
			return nil, &models.MetricError{
				Text: "header not allowed",
				Code: http.StatusBadRequest,
			}
		}
	}
	// validate params
	var metrics models.Metrics

	n := chi.URLParam(r, "name")
	t := chi.URLParam(r, "type")
	v := chi.URLParam(r, "value")

	metrics.ID = n
	metrics.MType = t

	if t == models.GaugeType {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, &models.MetricError{
				Text: "",
				Code: http.StatusBadRequest,
			}
		}
		metrics.Value = &value
	}

	if t == models.CounterType {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, &models.MetricError{
				Text: "",
				Code: http.StatusBadRequest,
			}
		}
		metrics.Delta = &value
	}

	if reflect.DeepEqual(metrics, models.Metrics{}) {
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&metrics); err != nil {
			logger.Log.Info("cannot decode request JSON body", zap.Error(err))
			return nil, &models.MetricError{
				Text: "cannot decode request JSON body",
				Code: http.StatusInternalServerError,
			}
		}
	}

	if metrics.ID == models.EmptyParam || metrics.MType == models.EmptyParam || ((metrics.MType == models.GaugeType && metrics.Value == nil) || (metrics.MType == models.CounterType && metrics.Delta == nil)) {
		return nil, &models.MetricError{
			Text: "empty param",
			Code: http.StatusBadRequest,
		}
	}

	if !slices.Contains([]string{models.GaugeType, models.CounterType}, metrics.MType) {
		return nil, &models.MetricError{
			Text: "incorrect type",
			Code: http.StatusBadRequest,
		}
	}

	return &metrics, nil
}

func sendResponse(enc *json.Encoder, data interface{}) {
	if err := enc.Encode(data); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
}
