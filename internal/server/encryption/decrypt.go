package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
	"os"
)

// DecryptSymmetric is a function that uses symmetric decryption to the given slice of bytes
func DecryptSymmetric(ciphertext []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
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

// DecryptAsymmetric is a function that uses asymmetric decryption to the given slice of bytes
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
