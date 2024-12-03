package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
)

// EncryptAsymmetric is a function that uses asymmetric encryption to the given slice of bytes
func EncryptAsymmetric(js []byte, path string) ([]byte, *domain.MetricsError) {

	publicKeyPEM, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Info("error encryption", zap.Error(err))
		return nil, &domain.MetricsError{Text: err.Error(), Code: http.StatusInternalServerError}
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, &domain.MetricsError{Text: err.Error(), Code: http.StatusInternalServerError}
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), js)
	if err != nil {
		return nil, &domain.MetricsError{Text: err.Error(), Code: http.StatusInternalServerError}
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

// GenKey generates crypto key
func GenKey(n int) ([]byte, error) {
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
