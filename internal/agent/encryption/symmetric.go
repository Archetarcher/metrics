package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
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

// Encrypt is a function that uses symmetric encryption to the given slice of bytes
func (s *Symmetric) Encrypt(plaintext []byte) []byte {
	block, err := aes.NewCipher([]byte(s.key))
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
