package encoding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"io"
	"net/http"
)

// RequestHashesMiddleware — middleware-hashes для входящих HTTP-запросов.
func RequestHashesMiddleware(next http.Handler, config *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hash := r.Header.Get("HashSHA256")

		if config.Key != domain.EmptyParam && hash != domain.EmptyParam {
			h := hmac.New(sha256.New, []byte(config.Key))
			body, err := io.ReadAll(r.Body)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			h.Write(body)
			sign := h.Sum(nil)

			s, err := hex.DecodeString(hash)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			if !hmac.Equal(s, sign) {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(rw, r.WithContext(r.Context()))

	})
}
