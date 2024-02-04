package mysql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

func TestNewCustomGormLogger(t *testing.T) {
	zapLog, _ := zap.NewDevelopment()
	l := NewCustomGormLogger(&options{
		requestIDKey: "request_id",
		gLog:         zapLog,
		logLevel:     logger.Info,
	})

	l.LogMode(logger.Info)
	ctx := context.WithValue(context.Background(), "request_id", "123")
	l.Info(ctx, "info", "foo")
	l.Warn(ctx, "warn", "bar")
	l.Error(ctx, "error", "foo bar")

	l.LogMode(logger.Silent)
	l.Trace(ctx, time.Now(), nil, nil)

	l.LogMode(logger.Info)
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "sql statement", 1
	}, nil)
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "sql statement", -1
	}, nil)

	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "sql statement", 0
	}, logger.ErrRecordNotFound)

	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "sql statement", 0
	}, errors.New("Error 1054: Unknown column 'test_column'"))

	l.LogMode(logger.Warn)
	l.Trace(ctx, time.Now(), func() (string, int64) {
		return "sql statement", 0
	}, logger.ErrRecordNotFound)
}

func Test_requestIDField(t *testing.T) {
	ctx := context.WithValue(context.Background(), "request_id", "123")
	field := requestIDField(ctx, "")
	assert.Equal(t, zap.Skip(), field)
	field = requestIDField(ctx, "your request id key")
	assert.Equal(t, zap.Skip(), field)
	field = requestIDField(ctx, "request_id")
	assert.Equal(t, zap.String("request_id", "123"), field)
}
