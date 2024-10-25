package encoding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/go-resty/resty/v2"

	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

func HashMiddleware(c *resty.Client, req *resty.Request, config *config.AppConfig) error {
	if config.Key != domain.EmptyParam && req.Header.Get("HashSHA256g") != domain.EmptyParam {

		h := hmac.New(sha256.New, []byte(config.Key))

		js, err := json.Marshal(req.Body)

		if err != nil {
			return err
		}

		h.Write(js)
		hash := h.Sum(nil)

		req.Header.Set(
			"HashSHA256g", hex.EncodeToString(hash))
	}

	return nil
}
