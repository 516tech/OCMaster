package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ShareCodeRepo struct {
	db *gorm.DB
}

func NewShareCodeRepo(db *gorm.DB) *ShareCodeRepo {
	return &ShareCodeRepo{db: db}
}

func (r *ShareCodeRepo) Exists(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&HardwareUploadModel{}).
		Where("share_code = ? AND expires_at > ?", code, time.Now()).Count(&count).Error
	return count > 0, err
}
