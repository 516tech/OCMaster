package mysql

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/ocmaster/backend/domain/suggestion"
)

type SuggestionRepo struct {
	db *gorm.DB
}

func NewSuggestionRepo(db *gorm.DB) *SuggestionRepo {
	return &SuggestionRepo{db: db}
}

func (r *SuggestionRepo) Save(ctx context.Context, s *suggestion.Suggestion) error {
	var cpuJSON, ramJSON, stJSON []byte
	if s.CpuSuggestion != nil {
		cpuJSON, _ = json.Marshal(s.CpuSuggestion)
	}
	if s.RamSuggestion != nil {
		ramJSON, _ = json.Marshal(s.RamSuggestion)
	}
	if s.StabilityTest != nil {
		stJSON, _ = json.Marshal(s.StabilityTest)
	}
	model := &SuggestionModel{
		MerchantID:    s.MerchantID,
		HardwareID:    s.HardwareID,
		CpuSuggestion: string(cpuJSON),
		RamSuggestion: string(ramJSON),
		StabilityTest: string(stJSON),
		RiskWarning:   s.RiskWarning,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *SuggestionRepo) FindByID(ctx context.Context, id uint64) (*suggestion.Suggestion, error) {
	var m SuggestionModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toSuggestionEntity(&m), nil
}

func (r *SuggestionRepo) FindByMerchant(ctx context.Context, merchantID uint64, offset, limit int) ([]*suggestion.Suggestion, int64, error) {
	var models []SuggestionModel
	var total int64
	r.db.WithContext(ctx).Model(&SuggestionModel{}).Where("merchant_id = ?", merchantID).Count(&total)
	if err := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	result := make([]*suggestion.Suggestion, len(models))
	for i, m := range models {
		result[i] = toSuggestionEntity(&m)
	}
	return result, total, nil
}

func (r *SuggestionRepo) DeleteOldHistory(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	result := r.db.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&QueryHistoryModel{})
	return result.RowsAffected, result.Error
}

func toSuggestionEntity(m *SuggestionModel) *suggestion.Suggestion {
	s := &suggestion.Suggestion{
		ID:         m.ID,
		MerchantID: m.MerchantID,
		HardwareID: m.HardwareID,
		RiskWarning: m.RiskWarning,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
	if m.CpuSuggestion != "" {
		var cpu suggestion.CpuSuggestion
		json.Unmarshal([]byte(m.CpuSuggestion), &cpu)
		s.CpuSuggestion = &cpu
	}
	if m.RamSuggestion != "" {
		var ram suggestion.RamSuggestion
		json.Unmarshal([]byte(m.RamSuggestion), &ram)
		s.RamSuggestion = &ram
	}
	if m.StabilityTest != "" {
		var st suggestion.StabilityTest
		json.Unmarshal([]byte(m.StabilityTest), &st)
		s.StabilityTest = &st
	}
	return s
}
