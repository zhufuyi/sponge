// Package dlock provides distributed locking primitives, supports redis and etcd.
package dlock

import "context"

// Locker is the interface that wraps the basic locking operations.
type Locker interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	TryLock(ctx context.Context) (bool, error)
	Close() error
}
