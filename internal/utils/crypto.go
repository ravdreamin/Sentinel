package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("invalid OTP length")
	}

	const charset = "0123456789"
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}

	return string(code), nil
}
