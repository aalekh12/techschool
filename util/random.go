package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate Random integer
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generate Random string
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i <= n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// Generate Money
func GenerateMoney() int64 {
	return RandomInt(10, 1000)
}

// Generate user
func GenerateUser() string {
	return RandomString(6)
}

// GenerateCurrency
func GenerateCurrency() string {
	currency := []string{"USD", "INR", "UER"}
	n := len(currency)

	return currency[rand.Intn(n)]
}
