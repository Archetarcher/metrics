package encryption

import (
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

type AsymmetricEncryption interface {
	Encrypt(text []byte) ([]byte, *domain.MetricsError)
}
type Asymmetric struct {
	keyPath string
}

func NewAsymmetric(path string) *Asymmetric {
	return &Asymmetric{keyPath: path}
}

// Encrypt is a function that uses asymmetric encryption to the given slice of bytes
func (a *Asymmetric) Encrypt(js []byte) ([]byte, *domain.MetricsError) {
	publicKeyPEM, err := os.ReadFile(a.keyPath)
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
