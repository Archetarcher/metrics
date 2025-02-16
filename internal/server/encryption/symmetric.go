package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
)

// SymmetricEncryption is interface that defines symmetric encryption.
type SymmetricEncryption interface {
	Encrypt(text []byte) []byte
}

// Symmetric is struct for symmetric encryption.
type Symmetric struct {
	key string
}

// NewSymmetric creates instance of Symmetric, use key for encryption.
func NewSymmetric(key string) *Symmetric {
	return &Symmetric{key: key}
}

// Decrypt is a function that uses symmetric decryption to the given slice of bytes
func (s *Symmetric) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(s.key))
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil, err
	}

	plaintext, err := gcm.Open(nil, ciphertext[len(ciphertext)-12:], ciphertext[:len(ciphertext)-12], nil)
	if err != nil {
		logger.Log.Info("error decryption", zap.Error(err))
		return nil, err
	}

	return plaintext, nil
}
