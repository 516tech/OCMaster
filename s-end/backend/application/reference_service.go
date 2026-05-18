package application

import (
	"context"

	"github.com/ocmaster/backend/domain/reference"
)

type ReferenceService struct {
	repo reference.Repository
}

func NewReferenceService(r reference.Repository) *ReferenceService {
	return &ReferenceService{repo: r}
}

func (s *ReferenceService) GetByCategory(ctx context.Context, category reference.Category) ([]*reference.ReferenceData, error) {
	return s.repo.FindByCategory(ctx, category)
}
