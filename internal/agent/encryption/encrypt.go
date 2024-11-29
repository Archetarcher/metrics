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
	"time"
)

// StartSession generates 16 bit key, encrypts it asymmetrically and sends to server.
// Saves key value to the config
func StartSession(config *config.AppConfig, client *resty.Client, retryCount int) *domain.TrackingError {
	url := "http://" + config.ServerRunAddr + "/session/"

	key, gErr := genKey(16)
	if gErr != nil {
		return &domain.TrackingError{Code: http.StatusInternalServerError, Text: "failed to generate crypto key"}
	}
	encryptedKey, eErr := EncryptAsymmetric(key, config.PublicKeyPath)
	if eErr != nil {
		return eErr
	}

	res, err := client.
		R().
		SetBody(domain.SessionRequest{Key: encryptedKey}).
		Post(url)

	if err != nil && retryCount > 0 {
		time.Sleep(time.Duration(config.ReportInterval) * time.Second)
		return StartSession(config, client, retryCount-1)
	}

	if err != nil {
		return &domain.TrackingError{Text: fmt.Sprintf("client: could not create request: %s\n", err.Error()), Code: http.StatusInternalServerError}
	}

	if res.StatusCode() != http.StatusOK && retryCount > 0 {
		time.Sleep(time.Duration(config.ReportInterval) * time.Second)

		return StartSession(config, client, retryCount-1)
	}

	if res.StatusCode() != http.StatusOK {
		return &domain.TrackingError{Text: fmt.Sprintf("client: responded with error creating session: %s\n, %s, %s", err, url, string(key)), Code: res.StatusCode()}
	}

	config.Session.Key = string(key)

	return nil
}

// EncryptAsymmetric is a function that uses asymmetric encryption to the given slice of bytes
func EncryptAsymmetric(js []byte, path string) ([]byte, *domain.TrackingError) {

	publicKeyPEM, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Info("error encryption", zap.Error(err))
		return nil, &domain.TrackingError{Text: err.Error(), Code: http.StatusInternalServerError}
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, &domain.TrackingError{Text: err.Error(), Code: http.StatusInternalServerError}
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), js)
	if err != nil {
		return nil, &domain.TrackingError{Text: err.Error(), Code: http.StatusInternalServerError}
	}
	return ciphertext, nil
}

// EncryptSymmetric is a function that uses symmetric encryption to the given slice of bytes
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
