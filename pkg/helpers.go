package pkg

import (
	"crypto/rand"
	"encoding/hex"
)

func GeneratePrefixedString(prefix string) string {
	value := make([]byte, 16)
	rand.Read(value)

	return prefix + hex.EncodeToString(value)
}
