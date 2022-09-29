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

// Init 连接redis
// dsn 支持格式：
// 没有密码，没有db：localhost:6379
// 有密码，有db：<user>:<pass>@localhost:6379/2
func Init(dsn string, opts ...Option) (*redis.Client, error) {
	o := defaultOptions()
	o.apply(opts...)

	opt, err := getRedisOpt(dsn, o)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)

	if o.enableTrace { // 根据设置是否开启链路跟踪
		rdb.AddHook(redisotel.TracingHook{})
	}

	return rdb, nil
}

func getRedisOpt(dsn string, opts *options) (*redis.Options, error) {
	if len(dsn) > 8 {
		if !strings.Contains(dsn[len(dsn)-3:], "/") {
			dsn += "/0" // 默认使用db 0
		}

		if dsn[:8] != "redis://" {
			dsn = "redis://" + dsn
		}
	}

	redisOpts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	if opts.dialTimeout > 0 {
		redisOpts.DialTimeout = opts.dialTimeout
	}
	if opts.readTimeout > 0 {
		redisOpts.ReadTimeout = opts.readTimeout
	}
	if opts.writeTimeout > 0 {
		redisOpts.WriteTimeout = opts.writeTimeout
	}

	return redisOpts, nil
}

// Init2 连接redis
func Init2(addr string, password string, db int, opts ...Option) *redis.Client {
	o := defaultOptions()
	o.apply(opts...)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  o.dialTimeout,
		ReadTimeout:  o.readTimeout,
		WriteTimeout: o.writeTimeout,
	})

	if o.enableTrace { // 根据设置是否开启链路跟踪
		rdb.AddHook(redisotel.TracingHook{})
	}

	return rdb
}
