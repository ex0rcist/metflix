package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/ex0rcist/metflix/internal/entities"
)

// RSA key to encrypt data
type PublicKey *rsa.PublicKey

// RSA key to decrypt data
type PrivateKey *rsa.PrivateKey

// NewPrivateKey reads RSA public key from file
func NewPrivateKey(path entities.FilePath) (PrivateKey, error) {
	pemBlock, err := readKey(path)
	if err != nil {
		return nil, err
	}

	var rsaKey PrivateKey
	switch pemBlock.Type {
	case "RSA PRIVATE KEY":
		rsaKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("security.NewPrivateKey - x509.ParsePKCS1PrivateKey: %w", err)

		}
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("security.NewPrivateKey - x509.ParsePKCS8PrivateKey: %w", err)
		}

		key, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("security.NewPrivateKey - key.(*rsa.PrivateKey): %w", err)
		}

		rsaKey = key.(*rsa.PrivateKey)
	default:
		return nil, fmt.Errorf("unknown key type %s", pemBlock.Type)
	}

	return rsaKey, nil
}

// NewPublicKey reads RSA public key from file
func NewPublicKey(path entities.FilePath) (PublicKey, error) {
	block, err := readKey(path)
	if err != nil {
		return nil, err
	}

	rawKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("security.NewPublicKey - x509.ParsePKIXPublicKey: %w", err)
	}

	key, ok := rawKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("security.NewPublicKey - .(*rsa.PublicKey): %w", entities.ErrBadRSAKey)
	}

	return key, nil
}

// Encrypt message with RSA using PublicKey
func Encrypt(src io.Reader, key PublicKey) (*bytes.Buffer, error) {
	msg := new(bytes.Buffer)

	chunkSize := (*rsa.PublicKey)(key).Size() - 2*sha256.New().Size() - 2
	chunk := make([]byte, chunkSize)

	for {
		n, err := src.Read(chunk)

		if n > 0 {
			// chop trailing zeroes
			if n != len(chunk) {
				chunk = chunk[:n]
			}

			encryptedChunk, encErr := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, chunk, nil)
			if encErr != nil {
				return nil, fmt.Errorf("security.Encrypt - rsa.EncryptOAEP: %w", encErr)
			}

			msg.Write(encryptedChunk)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("security.Encrypt - reader.Read: %w", err)
		}
	}

	return msg, nil
}

// Decrypt RSA-encoded message using PrivateKey
func Decrypt(src io.Reader, key PrivateKey) (*bytes.Buffer, error) {
	msg := new(bytes.Buffer)

	chunkSize := key.PublicKey.Size()
	chunk := make([]byte, chunkSize)

	for {
		n, err := src.Read(chunk)

		if n > 0 {
			// chop trailing zeroes
			if n != len(chunk) {
				chunk = chunk[:n]
			}

			decryptedChunk, decErr := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, chunk, nil)
			if decErr != nil {
				return nil, fmt.Errorf("security.Decrypt - rsa.DecryptOAEP: %w", decErr)
			}

			msg.Write(decryptedChunk)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("security.Decrypt - src.Read: %w", err)
		}
	}

	return msg, nil
}

// Read key in PEM format from file
func readKey(path entities.FilePath) (*pem.Block, error) {
	rawKey, err := os.ReadFile(path.String())
	if err != nil {
		return nil, fmt.Errorf("security.readKey - os.ReadFile: %w", err)
	}

	key, _ := pem.Decode(rawKey)
	if key == nil {
		return nil, fmt.Errorf("security.readKey - pem.Decode: %w", entities.ErrBadRSAKey)
	}

	return key, nil
}
