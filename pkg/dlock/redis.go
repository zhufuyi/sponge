package dlock

import (
	"context"
	"errors"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

// RedisLock implements Locker using Redis.
type RedisLock struct {
	mutex *redsync.Mutex
}

// NewRedisLock creates a new RedisLock.
func NewRedisLock(client *redis.Client, key string, options ...redsync.Option) (Locker, error) {
	if client == nil {
		return nil, errors.New("redis client is nil")
	}
	if key == "" {
		return nil, errors.New("key is empty")
	}
	return newLocker(client, key, options...), nil
}

// NewRedisClusterLock creates a new RedisClusterLock.
func NewRedisClusterLock(clusterClient *redis.ClusterClient, key string, options ...redsync.Option) (Locker, error) {
	if clusterClient == nil {
		return nil, errors.New("cluster redis client is nil")
	}
	if key == "" {
		return nil, errors.New("key is empty")
	}
	return newLocker(clusterClient, key, options...), nil
}

func newLocker(delegate redis.UniversalClient, key string, options ...redsync.Option) Locker {
	pool := goredis.NewPool(delegate)
	rs := redsync.New(pool)
	mutex := rs.NewMutex(key, options...)

	return &RedisLock{
		mutex: mutex,
	}
}

// TryLock tries to acquire the lock without blocking.
func (l *RedisLock) TryLock(ctx context.Context) (bool, error) {
	err := l.mutex.TryLockContext(ctx)
	if err == nil {
		return true, nil
	}
	return false, err
}

// Lock blocks until the lock is acquired or the context is canceled.
func (l *RedisLock) Lock(ctx context.Context) error {
	return l.mutex.LockContext(ctx)
}

// Unlock releases the lock, if unlocking the key is successful, the key will be automatically deleted
func (l *RedisLock) Unlock(ctx context.Context) error {
	_, err := l.mutex.UnlockContext(ctx)
	return err
}

// Close no-op for RedisLock.
func (l *RedisLock) Close() error {
	return nil
}
