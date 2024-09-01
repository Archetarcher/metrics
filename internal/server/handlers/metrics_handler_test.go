package handlers

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

var c = config.NewConfig(store.Config{Memory: &memory.Config{Active: true}, Pgx: &pgx.Config{}})

func TestMetricsHandler_UpdateMetrics(t *testing.T) {
	c.ParseConfig()
	ctx := context.Background()

	storage, err := store.NewStore(c.Store, ctx)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}

	repo := repositories.NewMetricsRepository(storage)
	service := services.NewMetricsService(repo)
	handler := NewMetricsHandler(service, c)
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)

	srv := httptest.NewServer(r)

	defer srv.Close()
	type request struct {
		query  string
		method string
		params map[string]string
	}
	type want struct {
		code int
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
				code: http.StatusOK,
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
				code: http.StatusOK,
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
				code: http.StatusBadRequest,
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
				code: http.StatusMethodNotAllowed,
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
				code: http.StatusNotFound,
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
				code: http.StatusNotFound,
			},
		},
		{
			name: "negative test invalid value  #7",
			request: request{
				query:  "/update/gauge/name/value",
				method: http.MethodPost,
				params: map[string]string{"type": "gauge", "name": "name", "value": "value"},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = srv.URL + tt.request.query
			req.SetPathParams(tt.request.params)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")
		})
	}
}
