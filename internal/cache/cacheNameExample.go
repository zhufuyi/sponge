package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/internal/model"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
)

// delete the templates code start
type keyTypeExample = string
type valueTypeExample = string

// delete the templates code end

const (
	// cache prefix key, must end with a colon
	cacheNameExampleCachePrefixKey = "prefixKeyExample:"
	// CacheNameExampleExpireTime expire time
	CacheNameExampleExpireTime = 10 * time.Minute // nolint
)

var _ CacheNameExampleCache = (*cacheNameExampleCache)(nil)

// CacheNameExampleCache cache interface
type CacheNameExampleCache interface { // nolint
	Set(ctx context.Context, keyNameExample keyTypeExample, valueNameExample valueTypeExample, duration time.Duration) error
	Get(ctx context.Context, keyNameExample keyTypeExample) (valueTypeExample, error)
	Del(ctx context.Context, keyNameExample keyTypeExample) error
}

type cacheNameExampleCache struct {
	cache cache.Cache
}

// NewCacheNameExampleCache create a new cache
func NewCacheNameExampleCache(cacheType *model.CacheType) CacheNameExampleCache {
	newObject := func() interface{} {
		return ""
	}
	cachePrefix := ""
	jsonEncoding := encoding.JSONEncoding{}

	var c cache.Cache
	if strings.ToLower(cacheType.CType) == "redis" {
		c = cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, newObject)
	} else {
		c = cache.NewMemoryCache(cachePrefix, jsonEncoding, newObject)
	}

	return &cacheNameExampleCache{
		cache: c,
	}
}

// cache key
func (c *cacheNameExampleCache) getCacheKey(keyNameExample keyTypeExample) string {
	return fmt.Sprintf("%s%v", cacheNameExampleCachePrefixKey, keyNameExample)
}

// Set cache
func (c *cacheNameExampleCache) Set(ctx context.Context, keyNameExample keyTypeExample, valueNameExample valueTypeExample, duration time.Duration) error {
	cacheKey := c.getCacheKey(keyNameExample)
	return c.cache.Set(ctx, cacheKey, &valueNameExample, duration)
}

// Get cache
func (c *cacheNameExampleCache) Get(ctx context.Context, keyNameExample keyTypeExample) (valueTypeExample, error) {
	var valueNameExample valueTypeExample
	cacheKey := c.getCacheKey(keyNameExample)
	err := c.cache.Get(ctx, cacheKey, &valueNameExample)
	if err != nil {
		return valueNameExample, err
	}
	return valueNameExample, nil
}

// Del delete cache
func (c *cacheNameExampleCache) Del(ctx context.Context, keyNameExample keyTypeExample) error {
	cacheKey := c.getCacheKey(keyNameExample)
	return c.cache.Del(ctx, cacheKey)
}
