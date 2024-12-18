package middlewares

import (
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}
	loggerResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// RequestLoggerMiddleware — log-middleware for incoming HTTP-requests.
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggerResponseWriter{
			ResponseWriter: rw,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r.WithContext(r.Context()))

		duration := time.Since(start)

		logger.Log.Info("got incoming HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),

			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)
	})
}

// Write writes to http.ResponseWriter.
func (r *loggerResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

// WriteHeader writes status to header.
func (r *loggerResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
