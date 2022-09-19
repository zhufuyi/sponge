package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Debug debug级别信息
func Debug(msg string, fields ...Field) {
	getLogger().Debug(msg, fields...)
}

// Info info级别信息
func Info(msg string, fields ...Field) {
	getLogger().Info(msg, fields...)
}

// Warn warn级别信息
func Warn(msg string, fields ...Field) {
	getLogger().Warn(msg, fields...)
}

// Error error级别信息
func Error(msg string, fields ...Field) {
	getLogger().Error(msg, fields...)
}

// Panic panic级别信息
func Panic(msg string, fields ...Field) {
	getLogger().Panic(msg, fields...)
}

// Fatal fatal级别信息
func Fatal(msg string, fields ...Field) {
	getLogger().Fatal(msg, fields...)
}

// Debugf 带格式化debug级别信息
func Debugf(format string, a ...interface{}) {
	getLogger().Debug(fmt.Sprintf(format, a...))
}

// Infof 带格式化info级别信息
func Infof(format string, a ...interface{}) {
	getLogger().Info(fmt.Sprintf(format, a...))
}

// Warnf 带格式化warn级别信息
func Warnf(format string, a ...interface{}) {
	getLogger().Warn(fmt.Sprintf(format, a...))
}

// Errorf 带格式化error级别信息
func Errorf(format string, a ...interface{}) {
	getLogger().Error(fmt.Sprintf(format, a...))
}

// Fatalf 带格式化fatal级别信息
func Fatalf(format string, a ...interface{}) {
	getLogger().Fatal(fmt.Sprintf(format, a...))
}

// WithFields 携带字段信息
func WithFields(fields ...Field) *zap.Logger {
	return getLogger().With(fields...)
}
