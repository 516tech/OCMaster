package mysql

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ocmaster/backend/domain/hardware"
	"github.com/ocmaster/backend/domain/merchant"
	"github.com/ocmaster/backend/domain/suggestion"
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

func TestMerchantRepo_FindByIDAndUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMerchantRepo(db)
	ctx := context.Background()

	m := &merchant.Merchant{
		Phone: "13800138000", PasswordHash: "hash", Status: merchant.StatusActive,
	}
	if err := repo.Save(ctx, m); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Save doesn't populate entity ID, so find by phone first
	saved, _ := repo.FindByPhone(ctx, "13800138000")
	found, err := repo.FindByID(ctx, saved.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.ID != saved.ID {
		t.Errorf("expected ID %d, got %d", saved.ID, found.ID)
	}

	found.RiskTemplate = "custom template"
	if err := repo.Update(ctx, found); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := repo.FindByID(ctx, saved.ID)
	if updated.RiskTemplate != "custom template" {
		t.Errorf("expected 'custom template', got %q", updated.RiskTemplate)
	}
}

func TestReferenceRepo_FindByCategoryAndKey(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReferenceRepo(db)
	ctx := context.Background()

	ref := &ReferenceDataModel{Category: "cpu", ModelKey: "i7-13700K", Params: `{"base_freq":"3.4GHz"}`}
	if err := db.WithContext(ctx).Create(ref).Error; err != nil {
		t.Fatalf("seed reference failed: %v", err)
	}

	items, err := repo.FindByCategory(ctx, "cpu")
	if err != nil {
		t.Fatalf("FindByCategory failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}

	found, err := repo.FindByKey(ctx, "cpu", "i7-13700K")
	if err != nil {
		t.Fatalf("FindByKey failed: %v", err)
	}
	if found.ModelKey != "i7-13700K" {
		t.Errorf("expected i7-13700K, got %s", found.ModelKey)
	}
}

func TestShareCodeRepo_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewShareCodeRepo(db)
	ctx := context.Background()

	hwRepo := NewHardwareRepo(db)
	hwRepo.Save(ctx, &hardware.HardwareUpload{
		ShareCode: "999888", HardwareData: hardware.HardwareInfo{}, ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	exists, err := repo.Exists(ctx, "999888")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("expected share code to exist")
	}

	exists2, _ := repo.Exists(ctx, "000000")
	if exists2 {
		t.Error("expected non-existent share code to return false")
	}
}

func TestSuggestionRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSuggestionRepo(db)
	ctx := context.Background()

	s := &suggestion.Suggestion{
		MerchantID: 1, HardwareID: 1,
		CpuSuggestion: &suggestion.CpuSuggestion{Frequency: "5.5GHz", Voltage: "1.35V"},
		RamSuggestion: &suggestion.RamSuggestion{Frequency: "3600MHz", Timings: "16-16-16-36"},
		StabilityTest: &suggestion.StabilityTest{Tool: "Prime95", Duration: "2h"},
		RiskWarning:   "高压注意",
	}
	if err := repo.Save(ctx, s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Save doesn't populate entity ID, find by merchant to get it
	items, _, _ := repo.FindByMerchant(ctx, 1, 0, 1)
	if len(items) == 0 {
		t.Fatal("expected at least 1 saved suggestion")
	}
	id := items[0].ID

	found, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.CpuSuggestion.Frequency != "5.5GHz" {
		t.Errorf("expected 5.5GHz, got %s", found.CpuSuggestion.Frequency)
	}
	if found.RamSuggestion.Frequency != "3600MHz" {
		t.Errorf("expected 3600MHz, got %s", found.RamSuggestion.Frequency)
	}
	if found.StabilityTest.Tool != "Prime95" {
		t.Errorf("expected Prime95, got %s", found.StabilityTest.Tool)
	}
}

func TestSuggestionRepo_FindByMerchant(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSuggestionRepo(db)
	ctx := context.Background()

	repo.Save(ctx, &suggestion.Suggestion{MerchantID: 1, RiskWarning: "a"})
	repo.Save(ctx, &suggestion.Suggestion{MerchantID: 1, RiskWarning: "b"})
	repo.Save(ctx, &suggestion.Suggestion{MerchantID: 2, RiskWarning: "c"})

	items, total, err := repo.FindByMerchant(ctx, 1, 0, 10)
	if err != nil {
		t.Fatalf("FindByMerchant failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestSuggestionRepo_DeleteOldHistory(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSuggestionRepo(db)
	ctx := context.Background()

	oldRecord := &QueryHistoryModel{CreatedAt: time.Now().Add(-40 * 24 * time.Hour)}
	if err := db.WithContext(ctx).Create(oldRecord).Error; err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	n, err := repo.DeleteOldHistory(ctx, 30)
	if err != nil {
		t.Fatalf("DeleteOldHistory failed: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 deleted, got %d", n)
	}
}
