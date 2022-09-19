package cache

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/dgraph-io/ristretto"
)

type memoryCache struct {
	client    *ristretto.Cache
	KeyPrefix string
	encoding  encoding.Encoding
}

// NewMemoryCache create a memory cache
func NewMemoryCache(keyPrefix string, encoding encoding.Encoding) Cache {
	// see: https://dgraph.io/blog/post/introducing-ristretto-high-perf-go-cache/
	//		https://www.start.io/blog/we-chose-ristretto-cache-for-go-heres-why/
	config := &ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	}
	store, _ := ristretto.NewCache(config)
	return &memoryCache{
		client:    store,
		KeyPrefix: keyPrefix,
		encoding:  encoding,
	}
}

// Set add cache
func (m *memoryCache) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	buf, err := encoding.Marshal(m.encoding, val)
	if err != nil {
		return fmt.Errorf("encoding.Marshal error: %v, key=%s, val=%+v ", err, key, val)
	}
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
	}
	m.client.SetWithTTL(cacheKey, buf, 0, expiration)
	return nil
}

// Get data
func (m *memoryCache) Get(ctx context.Context, key string, val interface{}) error {
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
	}
	data, ok := m.client.Get(cacheKey)
	if !ok {
		return nil
	}
	if data == NotFoundPlaceholder {
		return ErrPlaceholder
	}
	err = encoding.Unmarshal(m.encoding, data.([]byte), val)
	if err != nil {
		return fmt.Errorf("encoding.Unmarshal error: %v, key=%s, cacheKey=%s, type=%v, json=%+v ",
			err, key, cacheKey, reflect.TypeOf(val), string(data.([]byte)))
	}
	return nil
}

// Del 删除
func (m *memoryCache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	key := keys[0]
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		logger.Warn("build cache key error", logger.Err(err), logger.String("key", key))
		return err
	}
	m.client.Del(cacheKey)
	return nil
}

// MultiSet 批量set
func (m *memoryCache) MultiSet(ctx context.Context, valMap map[string]interface{}, expiration time.Duration) error {
	panic("implement me")
}

// MultiGet 批量获取
func (m *memoryCache) MultiGet(ctx context.Context, keys []string, val interface{}) error {
	panic("implement me")
}

func (m *memoryCache) SetCacheWithNotFound(ctx context.Context, key string) error {
	if m.client.Set(key, NotFoundPlaceholder, int64(DefaultNotFoundExpireTime)) {
		return nil
	}
	return ErrSetMemoryWithNotFound
}
