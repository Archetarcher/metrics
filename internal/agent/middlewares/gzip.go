package middlewares

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"github.com/go-resty/resty/v2"
)

// GzipMiddleware is a middleware for encrypting data before sending to server.
func GzipMiddleware(c *resty.Client, req *resty.Request, config *config.AppConfig) error {
	if req.Header.Get("Content-Encoding") != "gzip" {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)

		js, err := json.Marshal(req.Body)
		if err != nil {
			return err
		}

		_, err = zb.Write(js)
		if err != nil {
			return err
		}

		err = zb.Close()
		if err != nil {
			return err
		}

		compressed := buf.Bytes()
		req.Header.Set(
			"Content-Encoding", "gzip")

		//req.SetBody(compressed)

		// encryption
		encrypted := encryption.EncryptSymmetric(compressed, config.Session.Key)

		req.Header.Set(
			"Encrypted", "1")
		req.SetBody(encrypted)
	}

	return nil
}
