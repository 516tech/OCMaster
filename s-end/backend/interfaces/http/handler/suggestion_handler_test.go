package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/suggestion"
	"github.com/ocmaster/backend/infrastructure/auth"
	"github.com/ocmaster/backend/interfaces/http/middleware"
)

type mockSugRepo struct {
	saved *suggestion.Suggestion
	items []*suggestion.Suggestion
}

func (m *mockSugRepo) Save(ctx context.Context, s *suggestion.Suggestion) error { s.ID = 1; m.saved = s; return nil }
func (m *mockSugRepo) FindByID(ctx context.Context, id uint64) (*suggestion.Suggestion, error) {
	if id == 1 { return m.saved, nil }
	return nil, nil
}
func (m *mockSugRepo) FindByMerchant(ctx context.Context, mid uint64, o, l int) ([]*suggestion.Suggestion, int64, error) {
	return m.items, int64(len(m.items)), nil
}
func (m *mockSugRepo) DeleteOldHistory(ctx context.Context, d int) (int64, error) { return 0, nil }

func TestSuggestionHandler_Create(t *testing.T) {
	repo := &mockSugRepo{}
	svc := application.NewSuggestionService(repo)
	h := NewSuggestionHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"hardware_id": 1,
		"cpu_suggestion": map[string]string{"frequency": "5.5GHz", "voltage": "1.35V"},
		"risk_warning": "测试风险提示",
	})
	req := httptest.NewRequest("POST", "/api/v1/suggestions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.ClaimsKey, &auth.Claims{MerchantID: 1, Phone: "138"})
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSuggestionHandler_GetByID(t *testing.T) {
	repo := &mockSugRepo{}
	repo.Save(context.Background(), &suggestion.Suggestion{RiskWarning: "test"})
	svc := application.NewSuggestionService(repo)
	h := NewSuggestionHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/suggestions/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSuggestionHandler_ListHistory(t *testing.T) {
	repo := &mockSugRepo{
		items: []*suggestion.Suggestion{{RiskWarning: "a"}, {RiskWarning: "b"}},
	}
	svc := application.NewSuggestionService(repo)
	h := NewSuggestionHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/suggestions/history?offset=0&limit=10", nil)
	ctx := context.WithValue(req.Context(), middleware.ClaimsKey, &auth.Claims{MerchantID: 1, Phone: "138"})
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h.ListHistory(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
