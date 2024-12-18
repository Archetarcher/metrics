package handlers

import (
	"context"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
)

var conf Config

type Config struct {
	c      *config.AppConfig
	server *httptest.Server
	err    error
	once   sync.Once
}

func (c *Config) setConfig() {
	c.once.Do(func() {
		c.c = config.NewConfig()

		server, err := setupConfigServer()

		c.server = server
		c.err = err
	})
}
func setupConfigServer() (*httptest.Server, error) {
	ctx := context.Background()

	storage, err := store.NewStore(ctx, conf.c)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return nil, err.Err
	}
	conf.c.PrivateKeyPath = "../../../private.pem"

	repo := repositories.NewMetricsRepository(storage)
	service := services.NewMetricsService(repo)
	handler := NewMetricsHandler(service, conf.c)
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)
	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)
	r.Get("/value/{type}/{name}", handler.GetMetrics)
	r.Get("/", handler.GetMetricsPage)
	r.Get("/ping", handler.GetPing)

	r.Post("/update/", handler.UpdateMetricsJSON)
	r.Post("/updates/", handler.UpdatesMetrics)
	r.Post("/value/", handler.GetMetricsJSON)
	r.Post("/session/", handler.StartSession)

	srv := httptest.NewServer(r)

	return srv, nil
}

func init() {
	conf.setConfig()
}

var (
	counter = int64(2896127014)
	gauge   = 0.31167763133187076
	values  = [8]domain.Metrics{
		{
			ID:    "counter_value",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_2",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_2",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_3",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_3",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_4",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_4",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
	}
)

func TestMetricsHandler_UpdateMetrics(t *testing.T) {
	require.NoError(t, conf.err, "failed to init server", conf.server, conf.err)
	type request struct {
		params map[string]string
		query  string
		method string
	}
	type want struct {
		code int
	}
	tests := []struct {
		request request
		name    string
		want    want
	}{
		{
			name: "positive test counter #1",
			request: request{
				query:  "/update/{type}/{name}/{value}",
				method: http.MethodPost,
				params: map[string]string{"type": values[0].MType, "name": values[0].ID, "value": fmt.Sprintf("%d", *values[0].Delta)},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "positive test gauge #2",
			request: request{
				query:  "/update/{type}/{name}/{value}",
				method: http.MethodPost,
				params: map[string]string{"type": values[1].MType, "name": values[1].ID, "value": strconv.FormatFloat(*values[1].Value, 'f', -1, 64)},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "negative test invalid type #3",
			request: request{
				query:  "/update/{type}/{name}/{value}",
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
				query:  "/update/{type}/{name}/{value}",
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
				query:  "/update/{type}/",
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
				query:  "/update//",
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
				query:  "/update/{type}/{name}/{value}",
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
			req.URL = conf.server.URL + tt.request.query
			req.SetPathParams(tt.request.params)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", req.PathParams, req.URL)
		})
	}
}

func TestMetricsHandler_GetMetrics(t *testing.T) {
	require.NoError(t, conf.err, "failed to init server", conf.server, conf.err)
	type request struct {
		params map[string]string
		query  string
		method string
	}
	type want struct {
		code int
	}
	tests := []struct {
		request request
		name    string
		want    want
	}{
		{
			name: "positive test counter #1",
			request: request{
				query:  "/value/{type}/{name}",
				method: http.MethodGet,
				params: map[string]string{"type": values[0].MType, "name": values[0].ID},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "positive test gauge #2",
			request: request{
				query:  "/value/{type}/{name}",
				method: http.MethodGet,
				params: map[string]string{"type": values[1].MType, "name": values[1].ID},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "negative test invalid type #3",
			request: request{
				query:  "/value/{type}/{name}",
				method: http.MethodGet,
				params: map[string]string{"type": "gauged", "name": "value", "value": "1"},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test invalid method  #4",
			request: request{
				query:  "/value/{type}/{name}",
				method: http.MethodPost,
				params: map[string]string{"type": "gauge", "name": "value"},
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "negative test with invalid name  #5",
			request: request{
				query:  "/value/{type}/{name}",
				method: http.MethodGet,
				params: map[string]string{"type": "gauge", "name": "value"},
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name: "negative test without type  #6",
			request: request{
				query:  "/value//",
				method: http.MethodGet,
				params: map[string]string{},
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = conf.server.URL + tt.request.query
			req.SetPathParams(tt.request.params)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", req.PathParams, req.URL, resp.String())
		})
	}
}

func TestMetricsHandler_UpdateMetricsJSON(t *testing.T) {
	require.NoError(t, conf.err, "failed to init server", conf.server, conf.err)
	type request struct {
		query  string
		method string
		body   string
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
				query:  "/update/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\", \"delta\": %d}", values[2].ID, values[2].MType, *values[2].Delta),
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "positive test gauge #2",
			request: request{
				query:  "/update/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\", \"value\": %f}", values[3].ID, values[3].MType, *values[3].Value),
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "negative test invalid type #3",
			request: request{
				query:  "/update/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", "value", "gauged"),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test invalid method  #4",
			request: request{
				query:  "/update/",
				method: http.MethodGet,
				body:   "{}",
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "negative test without  name  #5",
			request: request{
				query:  "/update/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", "", "gauge"),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test without type  #6",
			request: request{
				query:  "/update/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", "test_name", ""),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test  #7",
			request: request{
				query:  "/update/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\", \"value\": %f", values[3].ID, values[3].MType, *values[3].Value),
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = conf.server.URL + tt.request.query
			req.SetBody(tt.request.body)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", req.PathParams, req.URL, resp.String(), tt.request.body)
		})
	}
}

func TestMetricsHandler_GetMetricsJSON(t *testing.T) {
	require.NoError(t, conf.err, "failed to init server", conf.server, conf.err)
	type request struct {
		body   string
		query  string
		method string
	}
	type want struct {
		code int
	}
	tests := []struct {
		request request
		name    string
		want    want
	}{
		{
			name: "positive test counter #1",
			request: request{
				query:  "/value/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", values[2].ID, values[2].MType),
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "positive test gauge #2",
			request: request{
				query:  "/value/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", values[3].ID, values[3].MType),
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "negative test invalid type #3",
			request: request{
				query:  "/value/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", "value", "gauged"),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test invalid method  #4",
			request: request{
				query:  "/value/",
				method: http.MethodGet,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", "value", "gauge"),
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "negative test with invalid name  #5",
			request: request{
				query:  "/value/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", "value", "gauge"),
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name: "negative test without type  #6",
			request: request{
				query:  "/value/",
				method: http.MethodPost,
				body:   "{}",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test #7",
			request: request{
				query:  "/value/",
				method: http.MethodPost,
				body:   fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"", "value", "gauge"),
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = conf.server.URL + tt.request.query
			req.SetBody(tt.request.body)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", req.PathParams, req.URL, resp.String())
		})
	}
}

func TestMetricsHandler_UpdatesMetrics(t *testing.T) {
	require.NoError(t, conf.err, "failed to init server", conf.server, conf.err)
	type request struct {
		query  string
		method string
		body   []domain.Metrics
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
			name: "positive test#1",
			request: request{
				query:  "/updates/",
				method: http.MethodPost,
				body:   []domain.Metrics{values[3], values[4], values[5], values[6]},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "negative test invalid type #3",
			request: request{
				query:  "/updates/",
				method: http.MethodPost,
				body: []domain.Metrics{{
					ID:    "value",
					MType: "gauged",
					Delta: nil,
					Value: nil,
				}},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test invalid method  #4",
			request: request{
				query:  "/updates/",
				method: http.MethodGet,
				body:   []domain.Metrics{},
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "negative test without  name  #5",
			request: request{
				query:  "/updates/",
				method: http.MethodPost,
				body: []domain.Metrics{{
					ID:    "",
					MType: "gauge",
					Delta: nil,
					Value: nil,
				}},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test without type  #6",
			request: request{
				query:  "/updates/",
				method: http.MethodPost,
				body: []domain.Metrics{{
					ID:    "test_name",
					MType: "",
					Delta: nil,
					Value: nil,
				}},
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
			req.URL = conf.server.URL + tt.request.query
			req.SetBody(tt.request.body)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", req.PathParams, req.URL, resp.String())
		})
	}
}

func TestMetricsHandler_GetMetricsPage(t *testing.T) {
	require.NoError(t, conf.err, "failed to init server", conf.server, conf.err)
	type request struct {
		query  string
		method string
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
			name: "positive test#1",
			request: request{
				query:  "/",
				method: http.MethodGet,
			},
			want: want{
				code: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = conf.server.URL + tt.request.query

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", req.PathParams, req.URL, resp.String())
		})
	}
}

func TestMetricsHandler_GetPing(t *testing.T) {
	require.NoError(t, conf.err, "failed to init server", conf.server, conf.err)
	type request struct {
		query  string
		method string
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
			name: "positive test#1",
			request: request{
				query:  "/ping",
				method: http.MethodGet,
			},
			want: want{
				code: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = conf.server.URL + tt.request.query

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", req.PathParams, req.URL, resp.String())
		})
	}
}

func TestMetricsHandler_StartSession(t *testing.T) {
	require.NoError(t, conf.err, "failed to init server", conf.server, conf.err)
	type request struct {
		params map[string][]byte
		query  string
		method string
	}
	type want struct {
		code int
	}

	key := "xsaxsaxsa"
	encryptedKey, eErr := encryption.EncryptAsymmetric([]byte(key), "../../../public.pem")
	require.Nil(t, eErr)
	tests := []struct {
		request request
		name    string
		want    want
	}{
		{
			name: "positive test  #1",
			request: request{
				query:  "/session/",
				method: http.MethodPost,
				params: map[string][]byte{"key": encryptedKey},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "negative test  #2",
			request: request{
				query:  "/session/",
				method: http.MethodPost,
				params: map[string][]byte{"key": []byte("wrong text")},
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name: "negative test  #3",
			request: request{
				query:  "/session/",
				method: http.MethodGet,
				params: map[string][]byte{"key": encryptedKey},
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.request.method
			req.URL = conf.server.URL + tt.request.query
			req.SetBody(tt.request.params)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, resp.StatusCode(), "Код ответа не совпадает с ожидаемым", req.PathParams, req.URL)
		})
	}
}
