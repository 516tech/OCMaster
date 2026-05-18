package application

import (
	"context"
	"testing"
	"time"

	"github.com/ocmaster/backend/domain/hardware"
	"github.com/ocmaster/backend/domain/suggestion"
)

type mockCleanupHardwareRepo struct {
	deleteExpiredCount int64
}

func (r *mockCleanupHardwareRepo) Save(ctx context.Context, h *hardware.HardwareUpload) error { return nil }
func (r *mockCleanupHardwareRepo) FindByShareCode(ctx context.Context, code string) (*hardware.HardwareUpload, error) {
	return nil, nil
}
func (r *mockCleanupHardwareRepo) DeleteByShareCode(ctx context.Context, code string) error { return nil }
func (r *mockCleanupHardwareRepo) DeleteExpired(ctx context.Context) (int64, error) {
	return r.deleteExpiredCount, nil
}

type mockCleanupSuggestionRepo struct {
	deleteOldCount int64
}

func (r *mockCleanupSuggestionRepo) Save(ctx context.Context, s *suggestion.Suggestion) error { return nil }
func (r *mockCleanupSuggestionRepo) FindByID(ctx context.Context, id uint64) (*suggestion.Suggestion, error) {
	return nil, nil
}
func (r *mockCleanupSuggestionRepo) FindByMerchant(ctx context.Context, mid uint64, o, l int) ([]*suggestion.Suggestion, int64, error) {
	return nil, 0, nil
}
func (r *mockCleanupSuggestionRepo) DeleteOldHistory(ctx context.Context, d int) (int64, error) {
	return r.deleteOldCount, nil
}

func TestCleanupService_StartAndStop(t *testing.T) {
	hr := &mockCleanupHardwareRepo{deleteExpiredCount: 0}
	sr := &mockCleanupSuggestionRepo{deleteOldCount: 0}
	svc := NewCleanupService(hr, sr)

	ctx, cancel := context.WithCancel(context.Background())
	svc.Start(ctx)

	// Let it run at least one tick
	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)
	// Service stops without panic
}

func TestCleanupService_WithExpiredData(t *testing.T) {
	hr := &mockCleanupHardwareRepo{deleteExpiredCount: 5}
	sr := &mockCleanupSuggestionRepo{deleteOldCount: 3}
	svc := NewCleanupService(hr, sr)

	ctx, cancel := context.WithCancel(context.Background())
	svc.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)
	// Cleanup runs without panic even with data
}
