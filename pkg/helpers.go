package pkg

import (
	"crypto/rand"
	"encoding/hex"
)

func GeneratePrefixedString(prefix string) string {
	value := make([]byte, 16)
	//TODO: может как-то поменять
	if _, err := rand.Read(value); err != nil {
		panic(err)
	}
	return prefix + hex.EncodeToString(value)
}
