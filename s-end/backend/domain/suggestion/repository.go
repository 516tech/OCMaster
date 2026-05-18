package suggestion

import "context"

type Repository interface {
	Save(ctx context.Context, s *Suggestion) error
	FindByID(ctx context.Context, id uint64) (*Suggestion, error)
	FindByMerchant(ctx context.Context, merchantID uint64, offset, limit int) ([]*Suggestion, int64, error)
	DeleteOldHistory(ctx context.Context, days int) (int64, error)
}
