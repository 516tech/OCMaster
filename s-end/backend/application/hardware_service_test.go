package application

import (
	"context"
	"testing"
	"time"

	"github.com/ocmaster/backend/domain/hardware"
	"github.com/ocmaster/backend/pkg/apperrors"
)

type mockHardwareRepo struct {
	saved      *hardware.HardwareUpload
	findResult *hardware.HardwareUpload
	findErr    error
}

func (m *mockHardwareRepo) Save(ctx context.Context, h *hardware.HardwareUpload) error {
	m.saved = h
	return nil
}
func (m *mockHardwareRepo) FindByShareCode(ctx context.Context, code string) (*hardware.HardwareUpload, error) {
	return m.findResult, m.findErr
}
func (m *mockHardwareRepo) DeleteByShareCode(ctx context.Context, code string) error { return nil }
func (m *mockHardwareRepo) DeleteExpired(ctx context.Context) (int64, error)         { return 0, nil }

type mockShareCodeRepo struct{ exists bool }

func (m *mockShareCodeRepo) Exists(ctx context.Context, code string) (bool, error) { return m.exists, nil }

func TestHardwareService_Upload(t *testing.T) {
	svc := NewHardwareService(&mockHardwareRepo{}, &mockShareCodeRepo{})
	info := &hardware.HardwareInfo{CPU: &hardware.CpuInfo{Model: "i7-13700K"}}
	code, err := svc.Upload(context.Background(), info)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("expected 6-digit code, got %s", code)
	}
}

func TestHardwareService_GetByCode_NotFound(t *testing.T) {
	repo := &mockHardwareRepo{findErr: apperrors.ErrNotFound}
	svc := NewHardwareService(repo, &mockShareCodeRepo{})
	_, err := svc.GetByCode(context.Background(), "000000")
	if err != apperrors.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestHardwareService_GetByCode_Expired(t *testing.T) {
	repo := &mockHardwareRepo{
		findResult: &hardware.HardwareUpload{
			ShareCode: "123456",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		},
	}
	svc := NewHardwareService(repo, &mockShareCodeRepo{})
	_, err := svc.GetByCode(context.Background(), "123456")
	if err != apperrors.ErrShareCodeExpired {
		t.Errorf("expected ErrShareCodeExpired, got %v", err)
	}
}

func TestHardwareService_Delete(t *testing.T) {
	repo := &mockHardwareRepo{
		findResult: &hardware.HardwareUpload{
			ShareCode: "123456",
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		},
	}
	svc := NewHardwareService(repo, &mockShareCodeRepo{})
	if err := svc.DeleteByCode(context.Background(), "123456"); err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}

func TestHardwareService_GetByCode_Success(t *testing.T) {
	repo := &mockHardwareRepo{
		findResult: &hardware.HardwareUpload{
			ShareCode: "123456",
			HardwareData: hardware.HardwareInfo{
				CPU: &hardware.CpuInfo{Model: "i7-13700K", Cores: 16},
			},
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		},
	}
	svc := NewHardwareService(repo, &mockShareCodeRepo{})
	result, err := svc.GetByCode(context.Background(), "123456")
	if err != nil {
		t.Fatalf("GetByCode failed: %v", err)
	}
	if result.HardwareData.CPU.Model != "i7-13700K" {
		t.Errorf("expected i7-13700K, got %s", result.HardwareData.CPU.Model)
	}
}
