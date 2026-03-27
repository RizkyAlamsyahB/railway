package sensitivedata

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const ciphertextPrefix = "enc:v1:"

// FieldCipher performs authenticated encryption for sensitive string fields.
type FieldCipher struct {
	aead cipher.AEAD
}

// NewFieldCipher creates a cipher from a 32-byte AES-256 key.
func NewFieldCipher(key []byte) (*FieldCipher, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("sensitive data encryption key must be exactly 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES-GCM: %w", err)
	}

	return &FieldCipher{aead: aead}, nil
}

// EncryptString encrypts a plaintext string. Empty values are preserved.
func (c *FieldCipher) EncryptString(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := c.aead.Seal(nil, nonce, []byte(plaintext), nil)
	payload := append(nonce, ciphertext...)
	return ciphertextPrefix + base64.StdEncoding.EncodeToString(payload), nil
}

// DecryptString decrypts a ciphertext string. Legacy plaintext values are returned as-is.
func (c *FieldCipher) DecryptString(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	if len(value) < len(ciphertextPrefix) || value[:len(ciphertextPrefix)] != ciphertextPrefix {
		return value, nil
	}

	payload, err := base64.StdEncoding.DecodeString(value[len(ciphertextPrefix):])
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted value: %w", err)
	}

	nonceSize := c.aead.NonceSize()
	if len(payload) < nonceSize {
		return "", fmt.Errorf("encrypted value payload is too short")
	}

	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt encrypted value: %w", err)
	}

	return string(plaintext), nil
}
