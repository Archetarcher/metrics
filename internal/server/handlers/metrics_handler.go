package handlers

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"net/http"
	"slices"
	"strconv"
)

type MetricsHandler struct {
	Storage *store.MemStorage
}

func (h *MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {

	// validate
	request, err := validateRequest(r)
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Text))
		return
	}
	metricsService := services.MetricsService{Storage: h.Storage}

	_, err = metricsService.Update(request)
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Text))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func validateRequest(r *http.Request) (*domain.UpdateRequest, *domain.ApplicationError) {
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
	n := r.PathValue("name")
	t := r.PathValue("type")
	v := r.PathValue("value")

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

	return &domain.UpdateRequest{
		Name:  n,
		Type:  t,
		Value: value,
	}, nil
}
