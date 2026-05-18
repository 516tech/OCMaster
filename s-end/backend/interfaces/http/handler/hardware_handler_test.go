package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/hardware"
	"github.com/ocmaster/backend/pkg/apperrors"
)

type testHardwareRepo struct {
	saved *hardware.HardwareUpload
	find  *hardware.HardwareUpload
	findErr error
}

func (r *testHardwareRepo) Save(ctx context.Context, h *hardware.HardwareUpload) error { r.saved = h; return nil }
func (r *testHardwareRepo) FindByShareCode(ctx context.Context, code string) (*hardware.HardwareUpload, error) {
	return r.find, r.findErr
}
func (r *testHardwareRepo) DeleteByShareCode(ctx context.Context, code string) error { return nil }
func (r *testHardwareRepo) DeleteExpired(ctx context.Context) (int64, error) { return 0, nil }

type testShareCodeRepo struct{}

func (r *testShareCodeRepo) Exists(ctx context.Context, code string) (bool, error) { return false, nil }

func TestHardwareHandler_Upload(t *testing.T) {
	repo := &testHardwareRepo{}
	svc := application.NewHardwareService(repo, &testShareCodeRepo{})
	h := NewHardwareHandler(svc)

	body, _ := json.Marshal(hardware.HardwareInfo{
		CPU: &hardware.CpuInfo{Model: "i7-13700K", Cores: 16},
	})
	req := httptest.NewRequest("POST", "/api/v1/hardware/upload", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["share_code"] == "" {
		t.Error("expected share_code in response")
	}
}

func TestHardwareHandler_DeleteByCode(t *testing.T) {
	repo := &testHardwareRepo{}
	svc := application.NewHardwareService(repo, &testShareCodeRepo{})
	h := NewHardwareHandler(svc)

	req := httptest.NewRequest("DELETE", "/api/v1/hardware/123456", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("code", "123456")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.DeleteByCode(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHardwareHandler_GetByCode_NotFound(t *testing.T) {
	repo := &testHardwareRepo{findErr: apperrors.ErrNotFound}
	svc := application.NewHardwareService(repo, &testShareCodeRepo{})
	h := NewHardwareHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/hardware/000000", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("code", "000000")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.GetByCode(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestHardwareHandler_GetByCode_Success(t *testing.T) {
	upload := &hardware.HardwareUpload{
		ShareCode: "123456",
		HardwareData: hardware.HardwareInfo{
			CPU: &hardware.CpuInfo{Model: "i7-13700K", Cores: 16, Threads: 24},
		},
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	repo := &testHardwareRepo{find: upload}
	svc := application.NewHardwareService(repo, &testShareCodeRepo{})
	h := NewHardwareHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/hardware/123456", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("code", "123456")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.GetByCode(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
