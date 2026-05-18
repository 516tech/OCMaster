package application

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ocmaster/backend/domain/merchant"
	"github.com/ocmaster/backend/infrastructure/auth"
	"github.com/ocmaster/backend/pkg/apperrors"
)

type MerchantService struct {
	repo      merchant.Repository
	jwtSecret string
}

func NewMerchantService(r merchant.Repository, jwtSecret string) *MerchantService {
	return &MerchantService{repo: r, jwtSecret: jwtSecret}
}

func (s *MerchantService) Register(ctx context.Context, phone, password string) error {
	_, err := s.repo.FindByPhone(ctx, phone)
	if err == nil {
		return apperrors.ErrConflict
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	m := &merchant.Merchant{
		Phone:        phone,
		PasswordHash: string(hash),
		Status:       merchant.StatusPending,
	}
	return s.repo.Save(ctx, m)
}

func (s *MerchantService) Login(ctx context.Context, phone, password string) (string, error) {
	m, err := s.repo.FindByPhone(ctx, phone)
	if err != nil || !m.IsActive() {
		return "", apperrors.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(m.PasswordHash), []byte(password)); err != nil {
		return "", apperrors.ErrUnauthorized
	}
	return auth.GenerateToken(s.jwtSecret, m.ID, m.Phone)
}

func (s *MerchantService) GetProfile(ctx context.Context, id uint64) (*merchant.Merchant, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return m, nil
}

func (s *MerchantService) UpdateProfile(ctx context.Context, m *merchant.Merchant) error {
	return s.repo.Update(ctx, m)
}

func (s *MerchantService) ChangePassword(ctx context.Context, id uint64, oldPass, newPass string) error {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(m.PasswordHash), []byte(oldPass)); err != nil {
		return apperrors.ErrUnauthorized
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	m.PasswordHash = string(hash)
	return s.repo.Update(ctx, m)
}
