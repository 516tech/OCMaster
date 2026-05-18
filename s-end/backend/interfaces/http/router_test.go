package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/hardware"
	"github.com/ocmaster/backend/interfaces/http/handler"
)

type dummyHWRepo struct{}
func (d *dummyHWRepo) Save(ctx context.Context, h *hardware.HardwareUpload) error { return nil }
func (d *dummyHWRepo) FindByShareCode(ctx context.Context, code string) (*hardware.HardwareUpload, error) {
	return nil, nil
}
func (d *dummyHWRepo) DeleteByShareCode(ctx context.Context, code string) error { return nil }
func (d *dummyHWRepo) DeleteExpired(ctx context.Context) (int64, error) { return 0, nil }

type dummySCRepo struct{}
func (d *dummySCRepo) Exists(ctx context.Context, code string) (bool, error) { return false, nil }

func makeRouter() chi.Router {
	hwSvc := application.NewHardwareService(&dummyHWRepo{}, &dummySCRepo{})
	return NewRouter(RouterDeps{
		HardwareHandler:   handler.NewHardwareHandler(hwSvc),
		MerchantHandler:   &handler.MerchantHandler{},
		SuggestionHandler: &handler.SuggestionHandler{},
		ReferenceHandler:  &handler.ReferenceHandler{},
		JWTSecret:         "test-secret",
	})
}

func TestNewRouter_HealthEndpoint(t *testing.T) {
	r := makeRouter()
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestNewRouter_ProfileRequiresAuth(t *testing.T) {
	r := makeRouter()
	req := httptest.NewRequest("GET", "/api/v1/merchant/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}
