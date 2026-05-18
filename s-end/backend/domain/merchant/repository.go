package merchant

import "context"

type Repository interface {
	Save(ctx context.Context, m *Merchant) error
	FindByPhone(ctx context.Context, phone string) (*Merchant, error)
	FindByID(ctx context.Context, id uint64) (*Merchant, error)
	Update(ctx context.Context, m *Merchant) error
}
