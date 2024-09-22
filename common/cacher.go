package common

import "context"

type Cacher interface {
	Set(ctx context.Context, keys []string, values []any) error
	Get(ctx context.Context, keys ...string) ([]any, error)
}

type CacherCloser interface {
	Cacher
	Closer
}
