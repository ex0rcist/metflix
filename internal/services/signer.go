package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ex0rcist/metflix/internal/entities"
)

var _ Signer = SignerService{}

// Interface to Signer service
type Signer interface {
	CalculateSignature(data []byte) (string, error)
	VerifySignature(data []byte, hash string) (bool, error)
}

// Signer service
type SignerService struct {
	secret []byte
}

// Signer constructor
func NewSignerService(secret entities.Secret) SignerService {
	return SignerService{secret: []byte(secret)}
}

// Calculate signature of []byte
func (s SignerService) CalculateSignature(data []byte) (string, error) {
	mac := hmac.New(sha256.New, s.secret)

	_, err := mac.Write(data)
	if err != nil {
		return "", err
	}
	digest := mac.Sum(nil)

	return hex.EncodeToString(digest), nil
}

// Verify signature, provided for []byte.
func (s SignerService) VerifySignature(data []byte, hash string) (bool, error) {
	if len(hash) == 0 {
		return false, entities.ErrNoSignature
	}

	expected, err := s.CalculateSignature(data)
	if err != nil {
		return false, err
	}

	return expected == hash, nil
}
