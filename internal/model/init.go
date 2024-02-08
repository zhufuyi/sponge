// Package model is the initial database driver and define the data structures corresponding to the tables.
package model

import (
	"strings"
	"sync"
	"time"

	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/pkg/ggorm"
	"github.com/zhufuyi/sponge/pkg/goredis"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// ErrCacheNotFound No hit cache
	ErrCacheNotFound = redis.Nil

	// ErrRecordNotFound no records found
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	db    *gorm.DB
	once1 sync.Once

	redisCli *redis.Client
	once2    sync.Once

	cacheType *CacheType
	once3     sync.Once
)

// CacheType cache type
type CacheType struct {
	CType string        // cache type  memory or redis
	Rdb   *redis.Client // if CType=redis, Rdb cannot be empty
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
		once3.Do(func() {
			InitCache(config.Get().App.CacheType)
		})
	}

	return cacheType
}

// InitRedis connect redis
func InitRedis() {
	opts := []goredis.Option{
		goredis.WithDialTimeout(time.Duration(config.Get().Redis.DialTimeout) * time.Second),
		goredis.WithReadTimeout(time.Duration(config.Get().Redis.ReadTimeout) * time.Second),
		goredis.WithWriteTimeout(time.Duration(config.Get().Redis.WriteTimeout) * time.Second),
	}
	if config.Get().App.EnableTrace {
		opts = append(opts, goredis.WithEnableTrace())
	}

	var err error
	redisCli, err = goredis.Init(config.Get().Redis.Dsn, opts...)
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}
}

// GetRedisCli get redis client
func GetRedisCli() *redis.Client {
	if redisCli == nil {
		once2.Do(func() {
			InitRedis()
		})
	}

	return redisCli
}

// CloseRedis close redis
func CloseRedis() error {
	if redisCli == nil {
		return nil
	}

	err := redisCli.Close()
	if err != nil && err.Error() != redis.ErrClosed.Error() {
		return err
	}

	return nil
}

// ------------------------------------------------------------------------------------------

// todo generate initialisation database code here
// delete the templates code start

// InitDB connect database
func InitDB() {
	switch strings.ToLower(config.Get().Database.Driver) {
	case "mysql", "tidb":
		InitMysql()
	case "postgres", "postgresql":
		InitPostgresql()
	default:
		panic("unsupported database driver: " + config.Get().Database.Driver)
	}
}

// GetDB get db
func GetDB() *gorm.DB {
	if db == nil {
		once1.Do(func() {
			switch strings.ToLower(config.Get().Database.Driver) {
			case "mysql", "tidb":
				InitMysql()
			case "postgres", "postgresql":
				InitPostgresql()
			default:
				panic("unsupported database driver: " + config.Get().Database.Driver)
			}
		})
	}

	return db
}

// InitMysql connect mysql
func InitMysql() {
	opts := []ggorm.Option{
		ggorm.WithMaxIdleConns(config.Get().Database.Mysql.MaxIdleConns),
		ggorm.WithMaxOpenConns(config.Get().Database.Mysql.MaxOpenConns),
		ggorm.WithConnMaxLifetime(time.Duration(config.Get().Database.Mysql.ConnMaxLifetime) * time.Minute),
	}
	if config.Get().Database.Mysql.EnableLog {
		opts = append(opts,
			ggorm.WithLogging(logger.Get()),
			ggorm.WithLogRequestIDKey("request_id"),
		)
	}

	if config.Get().App.EnableTrace {
		opts = append(opts, ggorm.WithEnableTrace())
	}

	// setting mysql slave and master dsn addresses,
	// if there is no read/write separation, you can comment out the following piece of code
	opts = append(opts, ggorm.WithRWSeparation(
		config.Get().Database.Mysql.SlavesDsn,
		config.Get().Database.Mysql.MastersDsn...,
	))

	// add custom gorm plugin
	//opts = append(opts, ggorm.WithGormPlugin(yourPlugin))

	var dsn = utils.AdaptiveMysqlDsn(config.Get().Database.Mysql.Dsn)
	var err error
	db, err = ggorm.InitMysql(dsn, opts...)
	if err != nil {
		panic("ggorm.InitMysql error: " + err.Error())
	}
}

// InitPostgresql connect postgresql
func InitPostgresql() {
	opts := []ggorm.Option{
		ggorm.WithMaxIdleConns(config.Get().Database.Postgresql.MaxIdleConns),
		ggorm.WithMaxOpenConns(config.Get().Database.Postgresql.MaxOpenConns),
		ggorm.WithConnMaxLifetime(time.Duration(config.Get().Database.Postgresql.ConnMaxLifetime) * time.Minute),
	}
	if config.Get().Database.Postgresql.EnableLog {
		opts = append(opts,
			ggorm.WithLogging(logger.Get()),
			ggorm.WithLogRequestIDKey("request_id"),
		)
	}

	if config.Get().App.EnableTrace {
		opts = append(opts, ggorm.WithEnableTrace())
	}

	// add custom gorm plugin
	//opts = append(opts, ggorm.WithGormPlugin(yourPlugin))

	var dsn = utils.AdaptivePostgresqlDsn(config.Get().Database.Postgresql.Dsn)
	var err error
	db, err = ggorm.InitPostgresql(dsn, opts...)
	if err != nil {
		panic("ggorm.InitPostgresql error: " + err.Error())
	}
}

// delete the templates code end

// CloseDB close db
func CloseDB() error {
	return ggorm.CloseDB(db)
}
