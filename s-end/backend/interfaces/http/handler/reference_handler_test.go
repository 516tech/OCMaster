package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/reference"
)

type testRefRepo struct {
	data []*reference.ReferenceData
}

func (r *testRefRepo) FindByCategory(_ context.Context, cat reference.Category) ([]*reference.ReferenceData, error) {
	return r.data, nil
}
func (r *testRefRepo) FindByKey(_ context.Context, cat reference.Category, key string) (*reference.ReferenceData, error) {
	return nil, nil
}

func TestReferenceHandler_GetByCategory(t *testing.T) {
	svc := application.NewReferenceService(&testRefRepo{
		data: []*reference.ReferenceData{
			{ModelKey: "i7-13700K", Params: `{"freq":"5.5GHz","voltage":"1.35V"}`},
		},
	})
	h := NewReferenceHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/reference/cpu", nil)
	rec := httptest.NewRecorder()
	h.GetByCategory(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestReferenceHandler_GetByCategory_Empty(t *testing.T) {
	svc := application.NewReferenceService(&testRefRepo{})
	h := NewReferenceHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/reference/ram", nil)
	rec := httptest.NewRecorder()
	h.GetByCategory(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
