package reference

import "context"

type Repository interface {
	FindByCategory(ctx context.Context, category Category) ([]*ReferenceData, error)
	FindByKey(ctx context.Context, category Category, modelKey string) (*ReferenceData, error)
}
