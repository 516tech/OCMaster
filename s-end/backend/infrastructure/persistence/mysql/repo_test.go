package mysql

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ocmaster/backend/domain/hardware"
	"github.com/ocmaster/backend/domain/merchant"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := AutoMigrate(db); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestHardwareRepo_SaveAndFind(t *testing.T) {
	db := setupTestDB(t)
	repo := NewHardwareRepo(db)
	ctx := context.Background()

	upload := &hardware.HardwareUpload{
		ShareCode: "123456",
		HardwareData: hardware.HardwareInfo{
			CPU: &hardware.CpuInfo{Model: "i7-13700K", Cores: 16, Threads: 24},
		},
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := repo.Save(ctx, upload); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	found, err := repo.FindByShareCode(ctx, "123456")
	if err != nil {
		t.Fatalf("FindByShareCode failed: %v", err)
	}
	if found.ShareCode != "123456" {
		t.Errorf("expected 123456, got %s", found.ShareCode)
	}
	if found.HardwareData.CPU.Model != "i7-13700K" {
		t.Errorf("expected i7-13700K, got %s", found.HardwareData.CPU.Model)
	}
}

func TestHardwareRepo_DeleteExpired(t *testing.T) {
	db := setupTestDB(t)
	repo := NewHardwareRepo(db)
	ctx := context.Background()

	expired := &hardware.HardwareUpload{
		ShareCode: "111111",
		HardwareData: hardware.HardwareInfo{},
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	repo.Save(ctx, expired)

	valid := &hardware.HardwareUpload{
		ShareCode: "222222",
		HardwareData: hardware.HardwareInfo{},
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	repo.Save(ctx, valid)

	n, err := repo.DeleteExpired(ctx)
	if err != nil {
		t.Fatalf("DeleteExpired failed: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 deleted, got %d", n)
	}
}

func TestMerchantRepo_SaveAndFind(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMerchantRepo(db)
	ctx := context.Background()

	m := &merchant.Merchant{
		Phone:        "13800138000",
		PasswordHash: "$2a$10$hashed",
		Status:       merchant.StatusActive,
	}
	if err := repo.Save(ctx, m); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	found, err := repo.FindByPhone(ctx, "13800138000")
	if err != nil {
		t.Fatalf("FindByPhone failed: %v", err)
	}
	if !found.IsActive() {
		t.Error("expected active merchant")
	}
}
