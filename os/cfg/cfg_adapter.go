package cfg

import "context"

type Adapter interface {
	Available(ctx context.Context, resource ...string) (ok bool)

	Get(ctx context.Context, pattern string) (value interface{}, err error)

	Data(ctx context.Context) (data map[string]interface{}, err error)
}
