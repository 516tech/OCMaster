package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/merchant"
	"github.com/ocmaster/backend/infrastructure/auth"
	"github.com/ocmaster/backend/interfaces/http/middleware"
)

type mockMerchRepo struct {
	findByPhoneFn func(string) (*merchant.Merchant, error)
	findByIDFn    func(uint64) (*merchant.Merchant, error)
	saved         *merchant.Merchant
}

func (m *mockMerchRepo) Save(ctx context.Context, mc *merchant.Merchant) error { m.saved = mc; return nil }
func (m *mockMerchRepo) FindByPhone(ctx context.Context, phone string) (*merchant.Merchant, error) {
	if m.findByPhoneFn != nil { return m.findByPhoneFn(phone) }; return nil, nil
}
func (m *mockMerchRepo) FindByID(ctx context.Context, id uint64) (*merchant.Merchant, error) {
	if m.findByIDFn != nil { return m.findByIDFn(id) }; return nil, nil
}
func (m *mockMerchRepo) Update(ctx context.Context, mc *merchant.Merchant) error { return nil }

func TestMerchantHandler_Register(t *testing.T) {
	repo := &mockMerchRepo{findByPhoneFn: func(s string) (*merchant.Merchant, error) { return nil, errors.New("not found") }}
	svc := application.NewMerchantService(repo, "secret")
	h := NewMerchantHandler(svc)
	body, _ := json.Marshal(map[string]string{"phone": "13800138000", "password": "test123"})
	req := httptest.NewRequest("POST", "/api/v1/merchant/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Register(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestMerchantHandler_Login(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	repo := &mockMerchRepo{findByPhoneFn: func(s string) (*merchant.Merchant, error) {
		return &merchant.Merchant{ID: 1, Phone: "13800138000", PasswordHash: string(hash), Status: merchant.StatusActive}, nil
	}}
	svc := application.NewMerchantService(repo, "secret")
	h := NewMerchantHandler(svc)
	body, _ := json.Marshal(map[string]string{"phone": "13800138000", "password": "test123"})
	req := httptest.NewRequest("POST", "/api/v1/merchant/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Login(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestMerchantHandler_GetProfile(t *testing.T) {
	repo := &mockMerchRepo{findByIDFn: func(id uint64) (*merchant.Merchant, error) {
		return &merchant.Merchant{ID: 1, Phone: "13800138000", Status: merchant.StatusActive}, nil
	}}
	svc := application.NewMerchantService(repo, "secret")
	h := NewMerchantHandler(svc)

	token, _ := auth.GenerateToken("secret", 1, "13800138000")
	req := httptest.NewRequest("GET", "/api/v1/merchant/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	claims, _ := auth.ParseToken("secret", token)
	ctx := context.WithValue(req.Context(), middleware.ClaimsKey, claims)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h.GetProfile(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
