package contract

import (
	"context"
	"github.com/gohade/hade/framework"
	"time"
)

const CacheKey = "hade:cache"

type RememberFunc func(ctx context.Context, container framework.Container) (interface{}, error)

type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	GetObj(ctx context.Context, key string, model interface{}) error
	Many(ctx context.Context, keys []string) (map[string]string, error)

	Set(ctx context.Context, key string, val string, timeout time.Duration) error
	SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error
	SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error
	SetForever(ctx context.Context, key string, val string) error
	SetForeverObj(ctx context.Context, key string, val interface{}) error

	SetTTL(ctx context.Context, key string, timeout time.Duration) error
	GetTTL(ctx context.Context, key string) (time.Duration, error)

	Remember(ctx context.Context, key string, timeout time.Duration, rememberFunc RememberFunc, model interface{}) error

	Calc(ctx context.Context, key string, step int64) (int64, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)

	Delete(ctx context.Context, key string) error
	DeleteMany(ctx context.Context, keys []string) error
}
