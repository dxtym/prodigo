package utils_test

import (
	"prodigo/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomInt(t *testing.T) {
	got := utils.GenerateRandomInt(100)
	assert.Less(t, got, int64(100))
}

func TestGenerateRandomString(t *testing.T) {
	got := utils.GenerateRandomString(100)
	assert.Len(t, got, 100)
}

func BenchmarkGenerateRandomInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GenerateRandomInt(100)
	}
}

func BenchmarkGenerateRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GenerateRandomString(100)
	}
}
