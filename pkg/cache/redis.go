package cache

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/zhufuyi/sponge/pkg/encoding"
)

// CacheNotFound no hit cache
var CacheNotFound = redis.Nil

// redisCache redis cache object
type redisCache struct {
	client            *redis.Client
	KeyPrefix         string
	encoding          encoding.Encoding
	DefaultExpireTime time.Duration
	newObject         func() interface{}
}

// NewRedisCache new a cache, client parameter can be passed in for unit testing
func NewRedisCache(client *redis.Client, keyPrefix string, encode encoding.Encoding, newObject func() interface{}) Cache {
	return &redisCache{
		client:    client,
		KeyPrefix: keyPrefix,
		encoding:  encode,
		newObject: newObject,
	}
}

// Set one value
func (c *redisCache) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	buf, err := encoding.Marshal(c.encoding, val)
	if err != nil {
		return fmt.Errorf("encoding.Marshal error: %v, key=%s, val=%+v ", err, key, val)
	}

	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
	}
	//if expiration == 0 {
	//	expiration = DefaultExpireTime
	//}
	err = c.client.Set(ctx, cacheKey, buf, expiration).Err()
	if err != nil {
		return fmt.Errorf("c.client.Set error: %v, cacheKey=%s", err, cacheKey)
	}
	return nil
}

// Get one value
func (c *redisCache) Get(ctx context.Context, key string, val interface{}) error {
	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
	}

	bytes, err := c.client.Get(ctx, cacheKey).Bytes()
	// NOTE: don't handle the case where redis value is nil
	// but leave it to the upstream for processing
	if err != nil {
		return err
	}

	// prevent Unmarshal from reporting an error if data is empty
	if string(bytes) == "" {
		return nil
	}
	if string(bytes) == NotFoundPlaceholder {
		return ErrPlaceholder
	}
	err = encoding.Unmarshal(c.encoding, bytes, val)
	if err != nil {
		return fmt.Errorf("encoding.Unmarshal error: %v, key=%s, cacheKey=%s, type=%v, json=%+v ",
			err, key, cacheKey, reflect.TypeOf(val), string(bytes))
	}
	return nil
}

// MultiSet set multiple values
func (c *redisCache) MultiSet(ctx context.Context, valueMap map[string]interface{}, expiration time.Duration) error {
	if len(valueMap) == 0 {
		return nil
	}
	//if expiration == 0 {
	//	expiration = DefaultExpireTime
	//}

	// the key-value is paired and has twice the capacity of a map
	paris := make([]interface{}, 0, 2*len(valueMap))
	for key, value := range valueMap {
		buf, err := encoding.Marshal(c.encoding, value)
		if err != nil {
			fmt.Printf("encoding.Marshal error, %v, value:%v\n", err, value)
			continue
		}
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			fmt.Printf("BuildCacheKey error, %v, key:%v\n", err, key)
			continue
		}
		paris = append(paris, []byte(cacheKey))
		paris = append(paris, buf)
	}
	pipeline := c.client.Pipeline()
	err := pipeline.MSet(ctx, paris...).Err()
	if err != nil {
		return fmt.Errorf("pipeline.MSet error: %v", err)
	}
	for i := 0; i < len(paris); i = i + 2 {
		switch paris[i].(type) {
		case []byte:
			pipeline.Expire(ctx, string(paris[i].([]byte)), expiration)
		default:
			fmt.Printf("redis expire is unsupported key type: %+v\n", reflect.TypeOf(paris[i]))
		}
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline.Exec error: %v", err)
	}
	return nil
}

// MultiGet get multiple values
func (c *redisCache) MultiGet(ctx context.Context, keys []string, value interface{}) error {
	if len(keys) == 0 {
		return nil
	}
	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
		}
		cacheKeys[index] = cacheKey
	}
	values, err := c.client.MGet(ctx, cacheKeys...).Result()
	if err != nil {
		return fmt.Errorf("c.client.MGet error: %v, keys=%+v", err, cacheKeys)
	}

	// Injection into map via reflection
	valueMap := reflect.ValueOf(value)
	for i, v := range values {
		if v == nil {
			continue
		}
		object := c.newObject()
		err = encoding.Unmarshal(c.encoding, []byte(v.(string)), object)
		if err != nil {
			fmt.Printf("unmarshal data error: %+v, key=%s, cacheKey=%s type=%v\n", err, keys[i], cacheKeys[i], reflect.TypeOf(value))
			continue
		}
		valueMap.SetMapIndex(reflect.ValueOf(cacheKeys[i]), reflect.ValueOf(object))
	}
	return nil
}

// Del delete multiple values
func (c *redisCache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			continue
		}
		cacheKeys[index] = cacheKey
	}
	err := c.client.Del(ctx, cacheKeys...).Err()
	if err != nil {
		return fmt.Errorf("c.client.Del error: %v, keys=%+v", err, cacheKeys)
	}
	return nil
}

// SetCacheWithNotFound set value for notfound
func (c *redisCache) SetCacheWithNotFound(ctx context.Context, key string) error {
	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
	}

	return c.client.Set(ctx, cacheKey, NotFoundPlaceholder, DefaultNotFoundExpireTime).Err()
}

// BuildCacheKey construct a cache key with a prefix
func BuildCacheKey(keyPrefix string, key string) (string, error) {
	if key == "" {
		return "", errors.New("[cache] key should not be empty")
	}

	cacheKey := key
	if keyPrefix != "" {
		cacheKey = strings.Join([]string{keyPrefix, key}, ":")
	}

	return cacheKey, nil
}

// -------------------------------------------------------------------------------------------

// redisClusterCache redis cluster cache object
type redisClusterCache struct {
	client            *redis.ClusterClient
	KeyPrefix         string
	encoding          encoding.Encoding
	DefaultExpireTime time.Duration
	newObject         func() interface{}
}

// NewRedisClusterCache new a cache
func NewRedisClusterCache(client *redis.ClusterClient, keyPrefix string, encode encoding.Encoding, newObject func() interface{}) Cache {
	return &redisClusterCache{
		client:    client,
		KeyPrefix: keyPrefix,
		encoding:  encode,
		newObject: newObject,
	}
}

// Set one value
func (c *redisClusterCache) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	buf, err := encoding.Marshal(c.encoding, val)
	if err != nil {
		return fmt.Errorf("encoding.Marshal error: %v, key=%s, val=%+v ", err, key, val)
	}

	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
	}
	//if expiration == 0 {
	//	expiration = DefaultExpireTime
	//}
	err = c.client.Set(ctx, cacheKey, buf, expiration).Err()
	if err != nil {
		return fmt.Errorf("c.client.Set error: %v, cacheKey=%s", err, cacheKey)
	}
	return nil
}

// Get one value
func (c *redisClusterCache) Get(ctx context.Context, key string, val interface{}) error {
	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
	}

	bytes, err := c.client.Get(ctx, cacheKey).Bytes()
	// NOTE: don't handle the case where redis value is nil
	// but leave it to the upstream for processing
	if err != nil {
		return err
	}

	// prevent Unmarshal from reporting an error if data is empty
	if string(bytes) == "" {
		return nil
	}
	if string(bytes) == NotFoundPlaceholder {
		return ErrPlaceholder
	}
	err = encoding.Unmarshal(c.encoding, bytes, val)
	if err != nil {
		return fmt.Errorf("encoding.Unmarshal error: %v, key=%s, cacheKey=%s, type=%v, json=%+v ",
			err, key, cacheKey, reflect.TypeOf(val), string(bytes))
	}
	return nil
}

// MultiSet set multiple values
func (c *redisClusterCache) MultiSet(ctx context.Context, valueMap map[string]interface{}, expiration time.Duration) error {
	if len(valueMap) == 0 {
		return nil
	}

	// the key-value is paired and has twice the capacity of a map
	paris := make([]interface{}, 0, 2*len(valueMap))
	for key, value := range valueMap {
		buf, err := encoding.Marshal(c.encoding, value)
		if err != nil {
			fmt.Printf("encoding.Marshal error, %v, value:%v\n", err, value)
			continue
		}
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			fmt.Printf("BuildCacheKey error, %v, key:%v\n", err, key)
			continue
		}
		paris = append(paris, []byte(cacheKey))
		paris = append(paris, buf)
	}
	pipeline := c.client.Pipeline()
	err := pipeline.MSet(ctx, paris...).Err()
	if err != nil {
		return fmt.Errorf("pipeline.MSet error: %v", err)
	}
	for i := 0; i < len(paris); i = i + 2 {
		switch paris[i].(type) {
		case []byte:
			pipeline.Expire(ctx, string(paris[i].([]byte)), expiration)
		default:
			fmt.Printf("redis expire is unsupported key type: %+v\n", reflect.TypeOf(paris[i]))
		}
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline.Exec error: %v", err)
	}
	return nil
}

// MultiGet get multiple values
func (c *redisClusterCache) MultiGet(ctx context.Context, keys []string, value interface{}) error {
	if len(keys) == 0 {
		return nil
	}
	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
		}
		cacheKeys[index] = cacheKey
	}
	values, err := c.client.MGet(ctx, cacheKeys...).Result()
	if err != nil {
		return fmt.Errorf("c.client.MGet error: %v, keys=%+v", err, cacheKeys)
	}

	// Injection into map via reflection
	valueMap := reflect.ValueOf(value)
	for i, v := range values {
		if v == nil {
			continue
		}
		object := c.newObject()
		err = encoding.Unmarshal(c.encoding, []byte(v.(string)), object)
		if err != nil {
			fmt.Printf("unmarshal data error: %+v, key=%s, cacheKey=%s type=%v\n", err, keys[i], cacheKeys[i], reflect.TypeOf(value))
			continue
		}
		valueMap.SetMapIndex(reflect.ValueOf(cacheKeys[i]), reflect.ValueOf(object))
	}
	return nil
}

// Del delete multiple values
func (c *redisClusterCache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			continue
		}
		cacheKeys[index] = cacheKey
	}
	err := c.client.Del(ctx, cacheKeys...).Err()
	if err != nil {
		return fmt.Errorf("c.client.Del error: %v, keys=%+v", err, cacheKeys)
	}
	return nil
}

// SetCacheWithNotFound set value for notfound
func (c *redisClusterCache) SetCacheWithNotFound(ctx context.Context, key string) error {
	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return fmt.Errorf("BuildCacheKey error: %v, key=%s", err, key)
	}

	return c.client.Set(ctx, cacheKey, NotFoundPlaceholder, DefaultNotFoundExpireTime).Err()
}
