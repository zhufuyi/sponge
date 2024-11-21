package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/zhufuyi/sponge/internal/database"
	"github.com/zhufuyi/sponge/internal/model"
)

const (
	// cache prefix key, must end with a colon
	{{.TableNameCamelFCL}}CachePrefixKey = "{{.TableNameCamelFCL}}:"
	// {{.TableNameCamel}}ExpireTime expire time
	{{.TableNameCamel}}ExpireTime = 5 * time.Minute
)

var _ {{.TableNameCamel}}Cache = (*{{.TableNameCamelFCL}}Cache)(nil)

// {{.TableNameCamel}}Cache cache interface
type {{.TableNameCamel}}Cache interface {
	Set(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}, data *model.{{.TableNameCamel}}, duration time.Duration) error
	Get(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) (*model.{{.TableNameCamel}}, error)
	MultiGet(ctx context.Context, {{.ColumnNamePluralCamelFCL}} []{{.GoType}}) (map[{{.GoType}}]*model.{{.TableNameCamel}}, error)
	MultiSet(ctx context.Context, data []*model.{{.TableNameCamel}}, duration time.Duration) error
	Del(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) error
	SetPlaceholder(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) error
	IsPlaceholderErr(err error) bool
}

// {{.TableNameCamelFCL}}Cache define a cache struct
type {{.TableNameCamelFCL}}Cache struct {
	cache cache.Cache
}

// New{{.TableNameCamel}}Cache new a cache
func New{{.TableNameCamel}}Cache(cacheType *database.CacheType) {{.TableNameCamel}}Cache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.{{.TableNameCamel}}{}
		})
		return &{{.TableNameCamelFCL}}Cache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.{{.TableNameCamel}}{}
		})
		return &{{.TableNameCamelFCL}}Cache{cache: c}
	}

	return nil // no cache
}

// Get{{.TableNameCamel}}CacheKey cache key
func (c *{{.TableNameCamelFCL}}Cache) Get{{.TableNameCamel}}CacheKey({{.ColumnNameCamelFCL}} {{.GoType}}) string {
	{{if .IsStringType}}return {{.TableNameCamelFCL}}CachePrefixKey + {{.ColumnNameCamelFCL}}{{else}}return {{.TableNameCamelFCL}}CachePrefixKey + utils.{{.GoTypeFCU}}ToStr({{.ColumnNameCamelFCL}}){{end}}
}

// Set write to cache
func (c *{{.TableNameCamelFCL}}Cache) Set(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}, data *model.{{.TableNameCamel}}, duration time.Duration) error {
	if data == nil {
		return nil
	}
	cacheKey := c.Get{{.TableNameCamel}}CacheKey({{.ColumnNameCamelFCL}})
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *{{.TableNameCamelFCL}}Cache) Get(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) (*model.{{.TableNameCamel}}, error) {
	var data *model.{{.TableNameCamel}}
	cacheKey := c.Get{{.TableNameCamel}}CacheKey({{.ColumnNameCamelFCL}})
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *{{.TableNameCamelFCL}}Cache) MultiSet(ctx context.Context, data []*model.{{.TableNameCamel}}, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.Get{{.TableNameCamel}}CacheKey(v.{{.ColumnNameCamel}})
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is {{.ColumnNameCamelFCL}} value
func (c *{{.TableNameCamelFCL}}Cache) MultiGet(ctx context.Context, {{.ColumnNamePluralCamelFCL}} []{{.GoType}}) (map[{{.GoType}}]*model.{{.TableNameCamel}}, error) {
	var keys []string
	for _, v := range {{.ColumnNamePluralCamelFCL}} {
		cacheKey := c.Get{{.TableNameCamel}}CacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.{{.TableNameCamel}})
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[{{.GoType}}]*model.{{.TableNameCamel}})
	for _, {{.ColumnNameCamelFCL}} := range {{.ColumnNamePluralCamelFCL}} {
		val, ok := itemMap[c.Get{{.TableNameCamel}}CacheKey({{.ColumnNameCamelFCL}})]
		if ok {
			retMap[{{.ColumnNameCamelFCL}}] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *{{.TableNameCamelFCL}}Cache) Del(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) error {
	cacheKey := c.Get{{.TableNameCamel}}CacheKey({{.ColumnNameCamelFCL}})
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *{{.TableNameCamelFCL}}Cache) SetPlaceholder(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) error {
	cacheKey := c.Get{{.TableNameCamel}}CacheKey({{.ColumnNameCamelFCL}})
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *{{.TableNameCamelFCL}}Cache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
