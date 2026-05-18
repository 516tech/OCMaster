package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/ocmaster/backend/domain/merchant"
)

type MerchantRepo struct {
	db *gorm.DB
}

func NewMerchantRepo(db *gorm.DB) *MerchantRepo {
	return &MerchantRepo{db: db}
}

func (r *MerchantRepo) Save(ctx context.Context, m *merchant.Merchant) error {
	model := &MerchantModel{
		Phone:        m.Phone,
		PasswordHash: m.PasswordHash,
		Status:       string(m.Status),
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *MerchantRepo) FindByPhone(ctx context.Context, phone string) (*merchant.Merchant, error) {
	var m MerchantModel
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&m).Error; err != nil {
		return nil, err
	}
	return toMerchantEntity(&m), nil
}

func (r *MerchantRepo) FindByID(ctx context.Context, id uint64) (*merchant.Merchant, error) {
	var m MerchantModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toMerchantEntity(&m), nil
}

func (r *MerchantRepo) Update(ctx context.Context, m *merchant.Merchant) error {
	return r.db.WithContext(ctx).Model(&MerchantModel{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"phone":         m.Phone,
		"password_hash": m.PasswordHash,
		"risk_template": m.RiskTemplate,
		"status":        string(m.Status),
	}).Error
}

func toMerchantEntity(m *MerchantModel) *merchant.Merchant {
	return &merchant.Merchant{
		ID:           m.ID,
		Phone:        m.Phone,
		PasswordHash: m.PasswordHash,
		Status:       merchant.Status(m.Status),
		RiskTemplate: m.RiskTemplate,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
