package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field type
type Field = zapcore.Field

// Int type
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

// Int64 type
func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

// Uint type
func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

// Uint64 type
func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

// Uintptr type
func Uintptr(key string, val uintptr) Field {
	return zap.Uintptr(key, val)
}

// Float64 type
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

// Bool type
func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

// String type
func String(key string, val string) Field {
	return zap.String(key, val)
}

// Stringer type
func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

// Time type
func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

// Duration type
func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

// Err type
func Err(err error) Field {
	return zap.Error(err)
}

// Any type, if it is a composite type such as object, slice, map, etc., use Any
func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}
