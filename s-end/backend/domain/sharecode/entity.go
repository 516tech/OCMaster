package sharecode

import (
	"crypto/rand"
	"math/big"
	"time"
)

type ShareCode struct {
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func Generate() string {
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code[i] = byte('0') + byte(n.Int64())
	}
	return string(code)
}

func ExpiryDate() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}

func (s *ShareCode) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
