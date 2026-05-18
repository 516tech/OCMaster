package mysql

import (
	"time"

	"gorm.io/gorm"
)

type MerchantModel struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	Phone        string    `gorm:"column:phone;type:varchar(20);uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null"`
	Status       string    `gorm:"column:status;type:varchar(20);default:pending"`
	RiskTemplate string    `gorm:"column:risk_template;type:text"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (MerchantModel) TableName() string { return "merchants" }

type HardwareUploadModel struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	ShareCode    string    `gorm:"column:share_code;type:varchar(6);uniqueIndex;not null"`
	HardwareData string    `gorm:"column:hardware_data;type:json;not null"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (HardwareUploadModel) TableName() string { return "hardware_uploads" }

type SuggestionModel struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	MerchantID    uint64    `gorm:"column:merchant_id;not null;index"`
	HardwareID    uint64    `gorm:"column:hardware_id;not null"`
	CpuSuggestion string    `gorm:"column:cpu_suggestion;type:json"`
	RamSuggestion string    `gorm:"column:ram_suggestion;type:json"`
	StabilityTest string    `gorm:"column:stability_test;type:json"`
	RiskWarning   string    `gorm:"column:risk_warning;type:text"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SuggestionModel) TableName() string { return "suggestions" }

type QueryHistoryModel struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	MerchantID uint64    `gorm:"column:merchant_id;not null;index"`
	ShareCode  string    `gorm:"column:share_code;type:varchar(6);not null"`
	HardwareID uint64    `gorm:"column:hardware_id;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (QueryHistoryModel) TableName() string { return "query_history" }

type ReferenceDataModel struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	Category string `gorm:"column:category;type:varchar(20);not null;uniqueIndex:uk_category_model"`
	ModelKey string `gorm:"column:model_key;type:varchar(100);not null;uniqueIndex:uk_category_model"`
	Params   string `gorm:"column:params;type:json;not null"`
}

func (ReferenceDataModel) TableName() string { return "reference_data" }

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&MerchantModel{},
		&HardwareUploadModel{},
		&SuggestionModel{},
		&QueryHistoryModel{},
		&ReferenceDataModel{},
	)
}
