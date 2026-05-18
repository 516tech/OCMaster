package application

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/ocmaster/backend/domain/merchant"
)

type mockMerchantRepo struct {
	saveCalled  bool
	findByPhone func(string) (*merchant.Merchant, error)
	findByID    func(uint64) (*merchant.Merchant, error)
}

func (m *mockMerchantRepo) Save(ctx context.Context, mc *merchant.Merchant) error { m.saveCalled = true; return nil }
func (m *mockMerchantRepo) FindByPhone(ctx context.Context, phone string) (*merchant.Merchant, error) {
	if m.findByPhone != nil {
		return m.findByPhone(phone)
	}
	return nil, nil
}
func (m *mockMerchantRepo) FindByID(ctx context.Context, id uint64) (*merchant.Merchant, error) {
	if m.findByID != nil {
		return m.findByID(id)
	}
	return nil, nil
}
func (m *mockMerchantRepo) Update(ctx context.Context, mc *merchant.Merchant) error { return nil }

func TestMerchantService_Register(t *testing.T) {
	repo := &mockMerchantRepo{
		findByPhone: func(string) (*merchant.Merchant, error) { return nil, errors.New("not found") },
	}
	svc := NewMerchantService(repo, "secret")
	err := svc.Register(context.Background(), "13800138000", "password123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if !repo.saveCalled {
		t.Error("expected Save to be called")
	}
}

func TestMerchantService_Register_Duplicate(t *testing.T) {
	repo := &mockMerchantRepo{
		findByPhone: func(string) (*merchant.Merchant, error) {
			return &merchant.Merchant{Phone: "13800138000"}, nil
		},
	}
	svc := NewMerchantService(repo, "secret")
	err := svc.Register(context.Background(), "13800138000", "password123")
	if err == nil {
		t.Error("expected conflict error for duplicate registration")
	}
}

func TestMerchantService_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockMerchantRepo{
		findByPhone: func(string) (*merchant.Merchant, error) {
			return &merchant.Merchant{
				ID:           1,
				Phone:        "13800138000",
				PasswordHash: string(hash),
				Status:       merchant.StatusActive,
			}, nil
		},
	}
	svc := NewMerchantService(repo, "secret")
	token, err := svc.Login(context.Background(), "13800138000", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestMerchantService_GetProfile(t *testing.T) {
	repo := &mockMerchantRepo{
		findByID: func(id uint64) (*merchant.Merchant, error) {
			return &merchant.Merchant{ID: 1, Phone: "13800138000"}, nil
		},
	}
	svc := NewMerchantService(repo, "secret")
	m, err := svc.GetProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if m == nil {
		t.Error("expected non-nil merchant")
	}
}

func TestMerchantService_GetProfile_NotFound(t *testing.T) {
	repo := &mockMerchantRepo{}
	repo.findByID = func(id uint64) (*merchant.Merchant, error) {
		return nil, errors.New("not found")
	}
	svc := NewMerchantService(repo, "secret")
	_, err := svc.GetProfile(context.Background(), 999)
	if err == nil {
		t.Error("expected error for not found")
	}
}

func TestMerchantService_UpdateProfile(t *testing.T) {
	repo := &mockMerchantRepo{}
	svc := NewMerchantService(repo, "secret")
	err := svc.UpdateProfile(context.Background(), &merchant.Merchant{
		ID: 1, Phone: "13800138000", RiskTemplate: "new template",
	})
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}
}

func TestMerchantService_ChangePassword_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	repo := &mockMerchantRepo{
		findByID: func(id uint64) (*merchant.Merchant, error) {
			return &merchant.Merchant{ID: 1, PasswordHash: string(hash)}, nil
		},
	}
	svc := NewMerchantService(repo, "secret")
	err := svc.ChangePassword(context.Background(), 1, "oldpass", "newpass")
	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}
}

func TestMerchantService_ChangePassword_WrongOld(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	repo := &mockMerchantRepo{
		findByID: func(id uint64) (*merchant.Merchant, error) {
			return &merchant.Merchant{ID: 1, PasswordHash: string(hash)}, nil
		},
	}
	svc := NewMerchantService(repo, "secret")
	err := svc.ChangePassword(context.Background(), 1, "wrongpass", "newpass")
	if err == nil {
		t.Error("expected error for wrong old password")
	}
}

func TestMerchantService_Login_Fail(t *testing.T) {
	repo := &mockMerchantRepo{
		findByPhone: func(string) (*merchant.Merchant, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewMerchantService(repo, "secret")
	_, err := svc.Login(context.Background(), "13800138000", "password123")
	if err == nil {
		t.Error("expected error for login with non-existent user")
	}
}

var _ = merchant.StatusActive // ensure import used
