package utils

import (
	"math/rand"
	"strings"
)

func GenerateRandomInt(limit int64) int64 {
	return rand.Int63n(limit)
}

func GenerateRandomString(length int) string {
	const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var builder strings.Builder
	for range length {
		index := rand.Intn(len(alpha))
		builder.WriteByte(alpha[index])
	}

	return builder.String()
}
