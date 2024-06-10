package handlers

import (
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler_UpdateMetrics(t *testing.T) {

	type request struct {
		query  string
		method string
		params map[string]string
	}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "positive test counter #1",
			request: request{
				query:  "/update/counter/counters/1",
				method: http.MethodPost,
				params: map[string]string{"type": "counter", "name": "counters", "value": "1"},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "",
			},
		},
		{
			name: "positive test gauge #2",
			request: request{
				query:  "/update/gauge/value/1",
				method: http.MethodPost,
				params: map[string]string{"type": "gauge", "name": "value", "value": "1"},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "",
			},
		},
		{
			name: "negative test invalid type #3",
			request: request{
				query:  "/update/gauged/value/1",
				method: http.MethodPost,
				params: map[string]string{"type": "gauged", "name": "value", "value": "1"},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name: "negative test invalid method  #4",
			request: request{
				query:  "/update/gauge/value/1",
				method: http.MethodGet,
				params: map[string]string{"type": "gauge", "name": "value", "value": "1"},
			},
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name: "negative test without name  #5",
			request: request{
				query:  "/update/gauge/",
				method: http.MethodPost,
				params: map[string]string{"type": "gauge"},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name: "negative test without type  #6",
			request: request{
				query:  "/update/",
				method: http.MethodPost,
				params: map[string]string{},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name: "negative test invalid value  #6",
			request: request{
				query:  "/update/gauge/name/value",
				method: http.MethodPost,
				params: map[string]string{"type": "gauge", "name": "name", "value": "value"},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
	}
	repo := &repositories.MetricRepository{Storage: store.NewStorage()}
	service := &services.MetricsService{MetricRepositoryInterface: repo}
	handler := MetricsHandler{MetricsServiceInterface: service}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.request.method, "http://localhost:8080"+tt.request.query, nil)
			for key, value := range tt.request.params {
				r.SetPathValue(key, value)

			}
			w := httptest.NewRecorder()

			handler.UpdateMetrics(w, r)
			result := w.Result()
			assert.Equal(t, tt.want.code, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

		})
	}
}
