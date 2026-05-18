package application

import (
	"context"
	"errors"
	"testing"

	"github.com/ocmaster/backend/domain/suggestion"
)

type mockSugRepo struct {
	saved    *suggestion.Suggestion
	find     *suggestion.Suggestion
	findErr  error
	items    []*suggestion.Suggestion
	itemsTotal int64
}

func (m *mockSugRepo) Save(ctx context.Context, s *suggestion.Suggestion) error { m.saved = s; return nil }
func (m *mockSugRepo) FindByID(ctx context.Context, id uint64) (*suggestion.Suggestion, error) { return m.find, m.findErr }
func (m *mockSugRepo) FindByMerchant(ctx context.Context, mid uint64, o, l int) ([]*suggestion.Suggestion, int64, error) {
	if m.items != nil {
		return m.items, m.itemsTotal, nil
	}
	return nil, 0, nil
}
func (m *mockSugRepo) DeleteOldHistory(ctx context.Context, d int) (int64, error) { return 0, nil }

func TestSuggestionService_GetByID(t *testing.T) {
	target := &suggestion.Suggestion{RiskWarning: "found"}
	repo := &mockSugRepo{find: target}
	svc := NewSuggestionService(repo)
	sug, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if sug.RiskWarning != "found" {
		t.Errorf("expected 'found', got %q", sug.RiskWarning)
	}
}

func TestSuggestionService_GetByID_NotFound(t *testing.T) {
	repo := &mockSugRepo{findErr: errors.New("not found")}
	svc := NewSuggestionService(repo)
	_, err := svc.GetByID(context.Background(), 999)
	if err == nil {
		t.Error("expected error for not found")
	}
}

func TestSuggestionService_ListHistory(t *testing.T) {
	repo := &mockSugRepo{}
	repo.items = []*suggestion.Suggestion{{RiskWarning: "a"}, {RiskWarning: "b"}}
	repo.itemsTotal = 2
	svc := NewSuggestionService(repo)
	items, total, err := svc.ListHistory(context.Background(), 1, 0, 10)
	if err != nil {
		t.Fatalf("ListHistory failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

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
