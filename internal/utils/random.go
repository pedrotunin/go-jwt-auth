package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GetRandomString(numBytes int) (string, error) {
	res := make([]byte, numBytes)
	_, err := rand.Read(res)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(res), nil
}
