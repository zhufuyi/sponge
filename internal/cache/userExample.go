package cache

import (
	"context"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/internal/model"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/spf13/cast"
)

const (
	// PrefixUserExampleCacheKey cache prefix
	PrefixUserExampleCacheKey = "userExample:"
)

var _ UserExampleCache = (*userExampleCache)(nil)

// UserExampleCache cache interface
type UserExampleCache interface {
	Set(ctx context.Context, id uint64, data *model.UserExample, duration time.Duration) error
	Get(ctx context.Context, id uint64) (ret *model.UserExample, err error)
	MultiGet(ctx context.Context, ids []uint64) (map[string]*model.UserExample, error)
	MultiSet(ctx context.Context, data []*model.UserExample, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// userExampleCache define a cache struct
type userExampleCache struct {
	cache cache.Cache
}

// NewUserExampleCache new a cache
func NewUserExampleCache(cacheType *model.CacheType) UserExampleCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	var c cache.Cache
	if strings.ToLower(cacheType.CType) == "redis" {
		c = cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserExample{}
		})
	} else {
		c = cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserExample{}
		})
	}

	return &userExampleCache{
		cache: c,
	}
}

// GetUserExampleCacheKey cache key
func (c *userExampleCache) GetUserExampleCacheKey(id uint64) string {
	return PrefixUserExampleCacheKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *userExampleCache) Set(ctx context.Context, id uint64, data *model.UserExample, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserExampleCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *userExampleCache) Get(ctx context.Context, id uint64) (*model.UserExample, error) {
	var data *model.UserExample
	cacheKey := c.GetUserExampleCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *userExampleCache) MultiSet(ctx context.Context, data []*model.UserExample, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserExampleCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *userExampleCache) MultiGet(ctx context.Context, ids []uint64) (map[string]*model.UserExample, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserExampleCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.UserExample)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*model.UserExample)
	for _, v := range ids {
		val, ok := itemMap[c.GetUserExampleCacheKey(v)]
		if ok {
			retMap[cast.ToString(v)] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *userExampleCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserExampleCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *userExampleCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserExampleCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
