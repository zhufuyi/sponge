// Package goredis is a library wrapped on top of github.com/go-redis/redis.
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

// Init connecting to redis
// dsn supported formats.
// (1) no password, no db: localhost:6379
// (2) with password and db: <user>:<pass>@localhost:6379/2
// (3) redis://default:123456@localhost:6379/0?max_retries=3
// for more parameters see the redis source code for the setupConnParams function
func Init(dsn string, opts ...Option) (*redis.Client, error) {
	o := defaultOptions()
	o.apply(opts...)

	opt, err := getRedisOpt(dsn, o)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)

	if o.enableTrace {
		rdb.AddHook(redisotel.TracingHook{})
	}

	return rdb, nil
}

func getRedisOpt(dsn string, opts *options) (*redis.Options, error) {
	dsn = strings.ReplaceAll(dsn, " ", "")
	if len(dsn) > 8 {
		if !strings.Contains(dsn[len(dsn)-3:], "/") {
			dsn += "/0" // use db 0 by default
		}

		if dsn[:8] != "redis://" && dsn[:9] != "rediss://" {
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
	if opts.tlsConfig != nil {
		redisOpts.TLSConfig = opts.tlsConfig
	}

	return redisOpts, nil
}

// Init2 connecting to redis
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
		TLSConfig:    o.tlsConfig,
	})

	if o.enableTrace {
		rdb.AddHook(redisotel.TracingHook{})
	}

	return rdb
}

// InitSentinel connecting to redis for sentinel, all redis username and password are the same
func InitSentinel(masterName string, addrs []string, username string, password string, opts ...Option) *redis.Client {
	o := defaultOptions()
	o.apply(opts...)

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addrs,
		Username:      username,
		Password:      password,
		DialTimeout:   o.dialTimeout,
		ReadTimeout:   o.readTimeout,
		WriteTimeout:  o.writeTimeout,
		TLSConfig:     o.tlsConfig,
	})

	if o.enableTrace {
		rdb.AddHook(redisotel.TracingHook{})
	}

	return rdb
}

// InitCluster connecting to redis for cluster, all redis username and password are the same
func InitCluster(addrs []string, username string, password string, opts ...Option) *redis.ClusterClient {
	o := defaultOptions()
	o.apply(opts...)

	clusterRdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Username:     username,
		Password:     password,
		DialTimeout:  o.dialTimeout,
		ReadTimeout:  o.readTimeout,
		WriteTimeout: o.writeTimeout,
		TLSConfig:    o.tlsConfig,
	})

	if o.enableTrace {
		clusterRdb.AddHook(redisotel.TracingHook{})
	}

	return clusterRdb
}
