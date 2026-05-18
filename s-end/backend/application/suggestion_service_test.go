package application

import (
	"context"
	"testing"

	"github.com/ocmaster/backend/domain/suggestion"
)

type mockSugRepo struct {
	saved *suggestion.Suggestion
	find  *suggestion.Suggestion
	findErr error
}

func (m *mockSugRepo) Save(ctx context.Context, s *suggestion.Suggestion) error { m.saved = s; return nil }
func (m *mockSugRepo) FindByID(ctx context.Context, id uint64) (*suggestion.Suggestion, error) { return m.find, m.findErr }
func (m *mockSugRepo) FindByMerchant(ctx context.Context, mid uint64, o, l int) ([]*suggestion.Suggestion, int64, error) {
	return nil, 0, nil
}
func (m *mockSugRepo) DeleteOldHistory(ctx context.Context, d int) (int64, error) { return 0, nil }

func TestSuggestionService_Create(t *testing.T) {
	repo := &mockSugRepo{}
	svc := NewSuggestionService(repo)
	err := svc.Create(context.Background(), &suggestion.Suggestion{
		MerchantID: 1,
		HardwareID: 1,
		CpuSuggestion: &suggestion.CpuSuggestion{Frequency: "5.5GHz", Voltage: "1.35V"},
		RiskWarning: "请谨慎操作",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if repo.saved == nil {
		t.Error("expected Save to be called")
	}
}
