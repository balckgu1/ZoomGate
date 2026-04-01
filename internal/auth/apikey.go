package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const apiKeyPrefix = "sk-zg-"

// GenerateAPIKey creates a new API key and returns the full key and its SHA-256 hash.
func GenerateAPIKey() (fullKey string, hash string, prefix string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", "", fmt.Errorf("generate random bytes: %w", err)
	}
	raw := hex.EncodeToString(b)
	fullKey = apiKeyPrefix + raw
	hash = HashAPIKey(fullKey)
	prefix = fullKey[:12]
	return
}

// HashAPIKey returns the SHA-256 hash of an API key.
func HashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
