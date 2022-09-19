package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field 字段类型
type Field = zapcore.Field

// Int int类型
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

// Int64 int64类型
func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

// Uint uint类型
func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

// Uint64 uint64类型
func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

// Uintptr uintptr类型
func Uintptr(key string, val uintptr) Field {
	return zap.Uintptr(key, val)
}

// Float64 float64类型
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

// Bool bool类型
func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

// String string类型
func String(key string, val string) Field {
	return zap.String(key, val)
}

// Stringer stringer类型
func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

// Time time.Time类型
func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

// Duration time.Duration类型
func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

// Err err类型
func Err(err error) Field {
	return zap.Error(err)
}

// Any 任意类型，如果是对象、slice、map等复合类型，使用Any
func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}
