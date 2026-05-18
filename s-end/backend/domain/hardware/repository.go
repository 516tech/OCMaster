package hardware

import "context"

type Repository interface {
	Save(ctx context.Context, upload *HardwareUpload) error
	FindByShareCode(ctx context.Context, code string) (*HardwareUpload, error)
	DeleteByShareCode(ctx context.Context, code string) error
	DeleteExpired(ctx context.Context) (int64, error)
}
