package gen

import (
	"math/rand"
	"time"
)

const (
	alpha = "abcdefghjklmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Alpha generates a string using letters only
func Alpha() string {
	return randString(20, alpha)
}

func randString(n int, alphabet string) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = alphabet[rand.Intn(len(alpha))]
	}
	return string(bytes)
}
