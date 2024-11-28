package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
)

const emptyParam = ""

// RequestDecryptMiddleware — decryption middleware for incoming http requests.
func RequestDecryptMiddleware(next http.Handler, config *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		enc := r.Header.Get("Encrypted")
		rw.Header().Set("session", config.Session)

		if config.PrivateKeyPath != emptyParam && enc != emptyParam {

			c, err := io.ReadAll(r.Body)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			decrypted := DecryptSymmetric(c, config.Session)
			r.Body = io.NopCloser(bytes.NewReader(decrypted))
		}

		next.ServeHTTP(rw, r.WithContext(r.Context()))

	})
}

func DecryptSymmetric(ciphertext []byte, key string) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil
	}

	plaintext, err := gcm.Open(nil, ciphertext[len(ciphertext)-12:], ciphertext[:len(ciphertext)-12], nil)
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil
	}

	return plaintext
}

func DecryptAsymmetric(ciphertext []byte, path string) ([]byte, error) {
	privateKeyPEM, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil, err
	}
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil, err
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil, err
	}

	return decrypted, nil
}
