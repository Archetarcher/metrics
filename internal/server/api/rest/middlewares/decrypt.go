package middlewares

import (
	"bytes"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/encryption"
	"io"
	"net/http"
)

// RequestDecryptMiddleware — decryption middleware for incoming http requests.
func RequestDecryptMiddleware(next http.Handler, config *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		enc := r.Header.Get("Encrypted")

		if config.PrivateKeyPath != emptyParam && enc != emptyParam {

			c, err := io.ReadAll(r.Body)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			decrypted := encryption.DecryptSymmetric(c, config.Session)
			r.Body = io.NopCloser(bytes.NewReader(decrypted))
		}

		next.ServeHTTP(rw, r.WithContext(r.Context()))

	})
}
