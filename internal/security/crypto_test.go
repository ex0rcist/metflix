package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
)

func generateTestKeys() (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey = &privateKey.PublicKey
	return privateKey, publicKey, nil
}

func writePEMFile(path string, pemType string, keyBytes []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	return pem.Encode(file, &pem.Block{
		Type:  pemType,
		Bytes: keyBytes,
	})
}

func TestNewPrivateKeyAndNewPublicKey(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("Failed to generate test keys: %v", err)
	}

	// Write private key to temporary file
	privKeyFile, err := os.CreateTemp("", "private_key.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file for private key: %v", err)
	}
	defer func() {
		err := os.Remove(privKeyFile.Name())
		if err != nil {
			panic(err)
		}
	}()

	privKeyPEM := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := writePEMFile(privKeyFile.Name(), "RSA PRIVATE KEY", privKeyPEM); err != nil {
		t.Fatalf("Failed to write private key PEM: %v", err)
	}

	// Write public key to temporary file
	pubKeyFile, kerr := os.CreateTemp("", "public_key.pem")
	if kerr != nil {
		t.Fatalf("Failed to create temp file for public key: %v", kerr)
	}
	defer func() {
		err := os.Remove(pubKeyFile.Name())
		if err != nil {
			panic(err)
		}
	}()

	pubKeyBytes, merr := x509.MarshalPKIXPublicKey(publicKey)
	if merr != nil {
		t.Fatalf("Failed to marshal public key: %v", merr)
	}
	if werr := writePEMFile(pubKeyFile.Name(), "PUBLIC KEY", pubKeyBytes); werr != nil {
		t.Fatalf("Failed to write public key PEM: %v", werr)
	}

	// Test NewPrivateKey
	privKey, kerr := NewPrivateKey(entities.FilePath(privKeyFile.Name()))
	if kerr != nil {
		t.Fatalf("NewPrivateKey failed: %v", kerr)
	}
	if privKey == nil || privKey.N.Cmp(privateKey.N) != 0 {
		t.Fatalf("NewPrivateKey did not return the expected key")
	}

	// Test NewPublicKey
	pubKey, k2err := NewPublicKey(entities.FilePath(pubKeyFile.Name()))
	if k2err != nil {
		t.Fatalf("NewPublicKey failed: %v", k2err)
	}
	if pubKey == nil || pubKey.N.Cmp(publicKey.N) != 0 {
		t.Fatalf("NewPublicKey did not return the expected key")
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	// Generate test keys
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("Failed to generate test keys: %v", err)
	}

	// Sample message to encrypt
	message := []byte("This is a test message for encryption")

	// Encrypt the message
	src := bytes.NewReader(message)
	encrypted, err := Encrypt(src, publicKey)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt the message
	decrypted, err := Decrypt(encrypted, privateKey)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// Compare decrypted message with original
	if !bytes.Equal(decrypted.Bytes(), message) {
		t.Fatalf("Decrypted message does not match original: got %s, want %s", decrypted.Bytes(), message)
	}
}

func TestReadKey(t *testing.T) {
	// Generate test private key
	privateKey, _, err := generateTestKeys()
	if err != nil {
		t.Fatalf("Failed to generate test private key: %v", err)
	}

	// Write private key to a temporary file in PEM format
	privKeyFile, err := os.CreateTemp("", "private_key.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file for private key: %v", err)
	}
	defer func() {
		rErr := os.Remove(privKeyFile.Name())
		if rErr != nil {
			panic(rErr)
		}
	}()

	privKeyPEM := x509.MarshalPKCS1PrivateKey(privateKey)
	if merr := writePEMFile(privKeyFile.Name(), "RSA PRIVATE KEY", privKeyPEM); merr != nil {
		t.Fatalf("Failed to write private key PEM: %v", merr)
	}

	// Read the key back using readKey
	pemBlock, err := readKey(entities.FilePath(privKeyFile.Name()))
	if err != nil {
		t.Fatalf("readKey failed: %v", err)
	}

	// Verify the type and content of the PEM block
	if pemBlock.Type != "RSA PRIVATE KEY" {
		t.Fatalf("Expected PEM type 'RSA PRIVATE KEY', got: %v", pemBlock.Type)
	}
	if !bytes.Equal(pemBlock.Bytes, privKeyPEM) {
		t.Fatalf("PEM bytes do not match original private key bytes")
	}
}
