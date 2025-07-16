package tinval

import "context"

type ValidationClient interface {
	Validate(ctx context.Context, id TIN) error
}
