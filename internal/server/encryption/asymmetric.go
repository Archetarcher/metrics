package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
	"os"
)

// AsymmetricEncryption is an interface that defines asymmetric encryption.
type AsymmetricEncryption interface {
	Encrypt(text []byte) ([]byte, *domain.MetricsError)
}

// Asymmetric is a struct for asymmetric encryption.
type Asymmetric struct {
	keyPath string
}

// NewAsymmetric creates instance of Asymmetric.
func NewAsymmetric(path string) *Asymmetric {
	return &Asymmetric{keyPath: path}
}

// Decrypt is a function that uses asymmetric decryption to the given slice of bytes
func (a *Asymmetric) Decrypt(ciphertext []byte) ([]byte, error) {
	privateKeyPEM, err := os.ReadFile(a.keyPath)
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
