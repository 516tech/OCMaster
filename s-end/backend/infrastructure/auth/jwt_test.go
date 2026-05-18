package auth

import (
	"testing"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(secret, 1, "13800138000")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.MerchantID != 1 {
		t.Errorf("expected MerchantID 1, got %d", claims.MerchantID)
	}
	if claims.Phone != "13800138000" {
		t.Errorf("expected Phone 13800138000, got %s", claims.Phone)
	}
}

func TestParseToken_InvalidSecret(t *testing.T) {
	token, _ := GenerateToken("secret-a", 1, "138")
	_, err := ParseToken("secret-b", token)
	if err == nil {
		t.Error("expected error for invalid secret")
	}
}

func TestParseToken_Garbage(t *testing.T) {
	_, err := ParseToken("secret", "not.a.token")
	if err == nil {
		t.Error("expected error for garbage token")
	}
}
