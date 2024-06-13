package handlers

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
	"slices"
	"strconv"
)

type MetricsHandler struct {
	MetricsService
}

type MetricsService interface {
	Update(request *domain.MetricRequest) (*domain.MetricResponse, *domain.ApplicationError)
	GetValue(request *domain.MetricRequest) (*domain.MetricResponse, *domain.ApplicationError)
	GetAllValues() (string, *domain.ApplicationError)
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result.Value))

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

func validateGetRequest(r *http.Request) (*domain.MetricRequest, *domain.ApplicationError) {
	// validate method
	if r.Method != http.MethodGet {

		return nil, &domain.ApplicationError{
			Text: "method not allowed",
			Code: http.StatusMethodNotAllowed,
		}
	}

	// validate headers
	for k, v := range domain.AllowedHeaders {
		if h := r.Header.Get(k); h != v {
			return nil, &domain.ApplicationError{
				Text: "header not allowed",
				Code: http.StatusBadRequest,
			}
		}
	}

	// validate params
	n := chi.URLParam(r, "name")
	t := chi.URLParam(r, "type")

	if n == domain.EmptyParam || t == domain.EmptyParam {
		return nil, &domain.ApplicationError{
			Text: "empty param",
			Code: http.StatusBadRequest,
		}
	}

	if !slices.Contains([]string{domain.GaugeType, domain.CounterType}, t) {
		return nil, &domain.ApplicationError{
			Text: "incorrect type",
			Code: http.StatusBadRequest,
		}
	}

	return &domain.MetricRequest{
		Name: n,
		Type: t,
	}, nil
}
func validateRequest(r *http.Request) (*domain.MetricRequest, *domain.ApplicationError) {
	// validate method
	if r.Method != http.MethodPost {

		return nil, &domain.ApplicationError{
			Text: "method not allowed",
			Code: http.StatusMethodNotAllowed,
		}
	}

	// validate headers
	for k, v := range domain.AllowedHeaders {
		if h := r.Header.Get(k); h != v {
			return nil, &domain.ApplicationError{
				Text: "header not allowed",
				Code: http.StatusBadRequest,
			}
		}
	}

	// validate params
	n := chi.URLParam(r, "name")
	t := chi.URLParam(r, "type")
	v := chi.URLParam(r, "value")

	if n == domain.EmptyParam || t == domain.EmptyParam || v == domain.EmptyParam {
		return nil, &domain.ApplicationError{
			Text: "empty param",
			Code: http.StatusBadRequest,
		}
	}

	if !slices.Contains([]string{domain.GaugeType, domain.CounterType}, t) {
		return nil, &domain.ApplicationError{
			Text: "incorrect type",
			Code: http.StatusBadRequest,
		}
	}
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, &domain.ApplicationError{
			Text: "",
			Code: http.StatusBadRequest,
		}
	}

	return &domain.MetricRequest{
		Name:  n,
		Type:  t,
		Value: value,
	}, nil
}
