package goredis

import (
	"strings"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
)

const (
	// ErrRedisNotFound not exist in redis
	ErrRedisNotFound = redis.Nil
	// DefaultRedisName default redis name
	DefaultRedisName = "default"
)

// RedisClient redis 客户端
var RedisClient *redis.Client

// Init 连接redis
// redisURL 支持格式：
// 没有密码，没有db：localhost:6379
// 有密码，有db：<user>:<pass>@localhost:6379/2
func Init(redisURL string, opts ...Option) (*redis.Client, error) {
	o := defaultOptions()
	o.apply(opts...)

	if len(redisURL) > 8 {
		if !strings.Contains(redisURL[len(redisURL)-3:], "/") {
			redisURL += "/0" // 默认使用db 0
		}

		if redisURL[:8] != "redis://" {
			redisURL = "redis://" + redisURL
		}
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)

	if o.enableTrace { // 根据设置是否开启链路跟踪
		rdb.AddHook(redisotel.TracingHook{})
	}

	return rdb, nil
}
