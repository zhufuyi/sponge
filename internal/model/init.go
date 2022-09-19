package model

import (
	"sync"
	"time"

	"github.com/zhufuyi/sponge/config"
	"github.com/zhufuyi/sponge/pkg/goredis"
	"github.com/zhufuyi/sponge/pkg/mysql"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// ErrRecordNotFound 没有找到记录
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	db    *gorm.DB
	once1 sync.Once

	redisCli *redis.Client
	once2    sync.Once
)

// InitMysql 连接mysql
func InitMysql() {
	opts := []mysql.Option{
		mysql.WithSlowThreshold(time.Duration(config.Get().Mysql.SlowThreshold) * time.Millisecond),
		mysql.WithMaxIdleConns(config.Get().Mysql.MaxIdleConns),
		mysql.WithMaxOpenConns(config.Get().Mysql.MaxOpenConns),
		mysql.WithConnMaxLifetime(time.Duration(config.Get().Mysql.ConnMaxLifetime) * time.Minute),
	}
	if config.Get().Mysql.EnableLog {
		opts = append(opts, mysql.WithLog())
	}

	if config.Get().App.EnableTracing {
		opts = append(opts, mysql.WithEnableTrace())
	}

	var err error
	db, err = mysql.Init(config.Get().Mysql.Dsn, opts...)
	if err != nil {
		panic("mysql.Init error: " + err.Error())
	}
}

// GetDB 返回db对象
func GetDB() *gorm.DB {
	if db == nil {
		once1.Do(func() {
			InitMysql()
		})
	}

	return db
}

// CloseMysql 关闭mysql
func CloseMysql() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if sqlDB != nil {
			return sqlDB.Close()
		}
	}

	return nil
}

// InitRedis 连接redis
func InitRedis() {
	opts := []goredis.Option{}
	if config.Get().App.EnableTracing {
		opts = append(opts, goredis.WithEnableTrace())
	}

	var err error
	redisCli, err = goredis.Init(config.Get().Redis.Dsn, opts...)
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}
}

// GetRedisCli 返回redis client
func GetRedisCli() *redis.Client {
	if redisCli == nil {
		once2.Do(func() {
			InitRedis()
		})
	}

	return redisCli
}

// CloseRedis 关闭redis
func CloseRedis() error {
	err := redisCli.Close()
	if err != nil && err.Error() != redis.ErrClosed.Error() {
		return err
	}

	return nil
}
