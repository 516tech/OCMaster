package application

import (
	"context"
	"testing"

	"github.com/ocmaster/backend/domain/reference"
)

type mockRefRepo struct{ data []*reference.ReferenceData }

func (m *mockRefRepo) FindByCategory(ctx context.Context, c reference.Category) ([]*reference.ReferenceData, error) {
	return m.data, nil
}
func (m *mockRefRepo) FindByKey(ctx context.Context, c reference.Category, k string) (*reference.ReferenceData, error) {
	return nil, nil
}

func TestReferenceService_GetByCategory(t *testing.T) {
	repo := &mockRefRepo{data: []*reference.ReferenceData{
		{ID: 1, Category: reference.CategoryCPU, ModelKey: "i7-13700K", Params: `{"freq":"5.5GHz"}`},
	}}
	svc := NewReferenceService(repo)
	items, err := svc.GetByCategory(context.Background(), reference.CategoryCPU)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}
