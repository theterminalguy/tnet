package util

import (
	"crypto/rand"
	"encoding/hex"
)

// SecureRandomBytes returns securely generated random bytes.
func SecureRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
