package cache

import (
	"context"
	"errors"
	"time"
)

var (
	// DefaultExpireTime default expiry time
	DefaultExpireTime = time.Hour * 24
	// DefaultNotFoundExpireTime expiry time when result is empty 1 minute,
	// often used for cache time when data is empty (cache pass-through)
	DefaultNotFoundExpireTime = time.Minute
	// NotFoundPlaceholder placeholder
	NotFoundPlaceholder = "*"

	// DefaultClient generate a cache client, where keyPrefix is generally the business prefix
	DefaultClient Cache

	// ErrPlaceholder .
	ErrPlaceholder = errors.New("cache: placeholder")
	// ErrSetMemoryWithNotFound .
	ErrSetMemoryWithNotFound = errors.New("cache: set memory cache err for not found")
)

// Cache driver interface
type Cache interface {
	Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, val interface{}) error
	MultiSet(ctx context.Context, valMap map[string]interface{}, expiration time.Duration) error
	MultiGet(ctx context.Context, keys []string, valueMap interface{}) error
	Del(ctx context.Context, keys ...string) error
	SetCacheWithNotFound(ctx context.Context, key string) error
}

// Set data
func Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	return DefaultClient.Set(ctx, key, val, expiration)
}

// Get data
func Get(ctx context.Context, key string, val interface{}) error {
	return DefaultClient.Get(ctx, key, val)
}

// MultiSet multiple set data
func MultiSet(ctx context.Context, valMap map[string]interface{}, expiration time.Duration) error {
	return DefaultClient.MultiSet(ctx, valMap, expiration)
}

// MultiGet multiple get data
func MultiGet(ctx context.Context, keys []string, valueMap interface{}) error {
	return DefaultClient.MultiGet(ctx, keys, valueMap)
}

// Del multiple delete data
func Del(ctx context.Context, keys ...string) error {
	return DefaultClient.Del(ctx, keys...)
}

// SetCacheWithNotFound .
func SetCacheWithNotFound(ctx context.Context, key string) error {
	return DefaultClient.SetCacheWithNotFound(ctx, key)
}
