package rest

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"github.com/go-resty/resty/v2"
)

const emptyParam = ""

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
		encrypted := encryption.NewSymmetric(config.Session.Key).Encrypt(compressed)

		req.Header.Set(
			"Encrypted", "1")
		req.SetBody(encrypted)
	}

	return nil
}

// HashMiddleware is a middleware for hashing data by HashSHA256g before sending to server.
func HashMiddleware(c *resty.Client, req *resty.Request, config *config.AppConfig) error {
	if config.Key != emptyParam && req.Header.Get("HashSHA256g") != emptyParam {

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
