package utils

import (
	"crypto/rand"
	"math/big"
)

// Base62 Alphabet for clean, URL-safe short codes.
const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateShortCode generates a cryptographically secure random alphanumeric code of the specified length.
func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}

	result := make([]byte, length)
	maxIndex := big.NewInt(int64(len(base62Chars)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return "", err
		}
		result[i] = base62Chars[num.Int64()]
	}

	return string(result), nil
}
