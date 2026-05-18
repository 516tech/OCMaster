package sharecode

import "context"

type Repository interface {
	Exists(ctx context.Context, code string) (bool, error)
}
