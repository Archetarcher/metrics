package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func StartSession(config *config.AppConfig, client *resty.Client) *domain.TrackingError {
	url := "http://" + config.ServerRunAddr + "/session/"

	key, err := genKey(16)
	if err != nil {
		return &domain.TrackingError{Code: http.StatusInternalServerError, Text: "failed to generate crypto key"}
	}
	encryptedKey := EncryptAsymmetric(key, config.PublicKeyPath)

	res, err := client.
		R().
		SetBody(domain.SessionRequest{Key: encryptedKey}).
		Post(url)
	if err != nil {
		return &domain.TrackingError{Text: fmt.Sprintf("client: could not create request: %s\n", err.Error()), Code: http.StatusInternalServerError}
	}

	if res.StatusCode() != http.StatusOK {
		return &domain.TrackingError{Text: fmt.Sprintf("client: responded with error: %s\n, %s", err, url), Code: res.StatusCode()}
	}
	config.Session = string(key)

	return nil
}

func genKey(n int) ([]byte, error) {
	rnd := make([]byte, n)
	nrnd, err := rand.Read(rnd)
	if err != nil {
		return nil, err
	} else if nrnd != n {
		return nil, fmt.Errorf(`nrnd %d != n %d`, nrnd, n)
	}
	for i := range rnd {
		rnd[i] = 'A' + rnd[i]%26
	}
	return rnd, nil
}

func EncryptAsymmetric(js []byte, path string) []byte {
	publicKeyPEM, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Info("error encryption", zap.Error(err))
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		logger.Log.Info("error encryption", zap.Error(err))
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), js)
	if err != nil {
		logger.Log.Info("error encryption", zap.Error(err))
	}

	return ciphertext
}

func EncryptSymmetric(plaintext []byte, key string) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.Log.Info("error encryption", zap.Error(err))
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Log.Info("error encryption", zap.Error(err))
		return nil
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		logger.Log.Info("error encryption", zap.Error(err))
		return nil
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	ciphertext = append(ciphertext, nonce...)
	return ciphertext
}
