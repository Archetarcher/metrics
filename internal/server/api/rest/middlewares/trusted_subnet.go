package middlewares

import (
	"net/http"

	"github.com/Archetarcher/metrics.git/internal/server/config"
)

// RequestTrustedSubnet — hash-middleware for incoming http requests.
func RequestTrustedSubnet(next http.Handler, config *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if config.TrustedSubnet != emptyParam {
			ipStr := r.Header.Get("X-Real-IP")

			if ipStr == emptyParam || ipStr != config.TrustedSubnet {
				rw.WriteHeader(http.StatusForbidden)
				return
			}

		}

		next.ServeHTTP(rw, r.WithContext(r.Context()))

	})
}
