package utils

import (
	"crypto/rand"
)

func GenerateCode() string {
	const length = 5
	const digits = "0123456789"

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "00000" // fallback (редко)
	}

	code := ""
	for i := 0; i < length; i++ {
		code += string(digits[int(b[i])%10])
	}

	return code
}
