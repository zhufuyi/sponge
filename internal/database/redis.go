package database

import (
	"sync"
	"time"

	"github.com/go-dev-frame/sponge/pkg/goredis"
	"github.com/go-dev-frame/sponge/pkg/tracer"

	"github.com/go-dev-frame/sponge/internal/config"
)

var (
	// ErrCacheNotFound No hit cache
	ErrCacheNotFound = goredis.ErrRedisNotFound
)

var (
	redisCli     *goredis.Client
	redisCliOnce sync.Once

	cacheType     *CacheType
	cacheTypeOnce sync.Once
)

// CacheType cache type
type CacheType struct {
	CType string          // cache type  memory or redis
	Rdb   *goredis.Client // if CType=redis, Rdb cannot be empty
}

// InitCache initial cache
func InitCache(cType string) {
	cacheType = &CacheType{
		CType: cType,
	}

	if cType == "redis" {
		cacheType.Rdb = GetRedisCli()
	}
}

// GetCacheType get cacheType
func GetCacheType() *CacheType {
	if cacheType == nil {
		cacheTypeOnce.Do(func() {
			InitCache(config.Get().App.CacheType)
		})
	}

	return cacheType
}

// InitRedis connect redis
func InitRedis() {
	redisCfg := config.Get().Redis
	opts := []goredis.Option{
		goredis.WithDialTimeout(time.Duration(redisCfg.DialTimeout) * time.Second),
		goredis.WithReadTimeout(time.Duration(redisCfg.ReadTimeout) * time.Second),
		goredis.WithWriteTimeout(time.Duration(redisCfg.WriteTimeout) * time.Second),
	}
	if config.Get().App.EnableTrace {
		opts = append(opts, goredis.WithTracing(tracer.GetProvider()))
	}

	var err error
	redisCli, err = goredis.Init(redisCfg.Dsn, opts...)
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}
}

// GetRedisCli get redis client
func GetRedisCli() *goredis.Client {
	if redisCli == nil {
		redisCliOnce.Do(func() {
			InitRedis()
		})
	}

	return redisCli
}

// CloseRedis close redis
func CloseRedis() error {
	return goredis.Close(redisCli)
}
