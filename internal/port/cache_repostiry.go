package port

import (
	"golang.org/x/net/context"
	"time"
)

type CacheRepository interface {
	Has(ctx context.Context, key string) bool
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expire time.Duration) error
	Del(ctx context.Context, key string) error
}
