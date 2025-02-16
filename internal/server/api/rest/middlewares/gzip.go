package middlewares

import (
	"compress/gzip"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

// compressWriter will implement http.ResponseWriter.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns http.ResponseWriter.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Header writes header to gzip.Writer.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader writes code to the header of compressWriter.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes gzip.Writer.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader implements interface io.ReadCloser.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads gzip.Reader.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes gzip.Reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// GzipMiddleware encodes and decodes body according to gzip encryption.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ow := rw

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(rw)
			ow = cw
			defer func() {
				err := cw.Close()
				if err != nil {
					logger.Log.Error("failed to close compress writer", zap.Error(err))

				}
			}()
			ow.Header().Set("Content-Encoding", "gzip")
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				logger.Log.Error("failed to create reader", zap.Error(err))

				rw.WriteHeader(http.StatusInternalServerError)

				return
			}
			r.Body = cr

			defer func() {
				cErr := cr.Close()
				if cErr != nil {
					logger.Log.Error("failed to close compress reader", zap.Error(err))
				}
			}()
		}
		next.ServeHTTP(ow, r.WithContext(r.Context()))

	})
}
