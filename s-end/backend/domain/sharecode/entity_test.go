package sharecode

import (
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	code := Generate()
	if len(code) != 6 {
		t.Errorf("expected 6-digit code, got %s", code)
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("expected numeric code, got %s", code)
		}
	}
}

func TestGenerate_Uniqueness(t *testing.T) {
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code := Generate()
		if codes[code] {
			t.Errorf("duplicate code generated: %s", code)
		}
		codes[code] = true
	}
}

func TestShareCode_IsExpired(t *testing.T) {
	s := &ShareCode{ExpiresAt: ExpiryDate()}
	if s.IsExpired() {
		t.Error("fresh share code should not be expired")
	}
}

func TestExpiryDate(t *testing.T) {
	exp := ExpiryDate()
	expected := time.Now().Add(7 * 24 * time.Hour)
	if exp.Sub(expected).Abs() > 5*time.Second {
		t.Error("expiry date should be 7 days from now")
	}
}
