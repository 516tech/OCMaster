package mysql

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/ocmaster/backend/domain/hardware"
)

type HardwareRepo struct {
	db *gorm.DB
}

func NewHardwareRepo(db *gorm.DB) *HardwareRepo {
	return &HardwareRepo{db: db}
}

func (r *HardwareRepo) Save(ctx context.Context, h *hardware.HardwareUpload) error {
	data, _ := json.Marshal(h.HardwareData)
	model := &HardwareUploadModel{
		ShareCode:    h.ShareCode,
		HardwareData: string(data),
		ExpiresAt:    h.ExpiresAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *HardwareRepo) FindByShareCode(ctx context.Context, code string) (*hardware.HardwareUpload, error) {
	var m HardwareUploadModel
	if err := r.db.WithContext(ctx).Where("share_code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	var info hardware.HardwareInfo
	json.Unmarshal([]byte(m.HardwareData), &info)
	return &hardware.HardwareUpload{
		ID:           m.ID,
		ShareCode:    m.ShareCode,
		HardwareData: info,
		ExpiresAt:    m.ExpiresAt,
		CreatedAt:    m.CreatedAt,
	}, nil
}

func (r *HardwareRepo) DeleteByShareCode(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).Where("share_code = ?", code).Delete(&HardwareUploadModel{}).Error
}

func (r *HardwareRepo) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&HardwareUploadModel{})
	return result.RowsAffected, result.Error
}
