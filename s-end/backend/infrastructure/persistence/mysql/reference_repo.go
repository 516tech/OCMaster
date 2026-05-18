package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/ocmaster/backend/domain/reference"
)

type ReferenceRepo struct {
	db *gorm.DB
}

func NewReferenceRepo(db *gorm.DB) *ReferenceRepo {
	return &ReferenceRepo{db: db}
}

func (r *ReferenceRepo) FindByCategory(ctx context.Context, category reference.Category) ([]*reference.ReferenceData, error) {
	var models []ReferenceDataModel
	if err := r.db.WithContext(ctx).Where("category = ?", string(category)).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*reference.ReferenceData, len(models))
	for i, m := range models {
		result[i] = &reference.ReferenceData{
			ID:       m.ID,
			Category: reference.Category(m.Category),
			ModelKey: m.ModelKey,
			Params:   m.Params,
		}
	}
	return result, nil
}

func (r *ReferenceRepo) FindByKey(ctx context.Context, category reference.Category, modelKey string) (*reference.ReferenceData, error) {
	var m ReferenceDataModel
	if err := r.db.WithContext(ctx).Where("category = ? AND model_key = ?", string(category), modelKey).First(&m).Error; err != nil {
		return nil, err
	}
	return &reference.ReferenceData{
		ID:       m.ID,
		Category: reference.Category(m.Category),
		ModelKey: m.ModelKey,
		Params:   m.Params,
	}, nil
}
