package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func Hash(data []byte, key []byte) (string, error) {
	h := hmac.New(sha256.New, key)
	if _, err := h.Write(data); err != nil {
		return "", err
	}
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
}
