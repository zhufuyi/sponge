// Package glog provides a gorm logger implementation based on zap.
package glog

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	gLog         *zap.Logger
	requestIDKey string
	logLevel     logger.LogLevel
}

// NewCustomGormLogger custom gorm logger
func NewCustomGormLogger(l *zap.Logger, requestIDKey string, logLevel logger.LogLevel) logger.Interface {
	if l == nil {
		l, _ = zap.NewProduction()
	}
	if requestIDKey == "" {
		requestIDKey = "request_id"
	}
	if logLevel == 0 {
		logLevel = logger.Info
	}
	return &gormLogger{
		gLog:         l,
		requestIDKey: requestIDKey,
		logLevel:     logLevel,
	}
}

// LogMode log mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logLevel = level
	return l
}

// Info print info
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		msg = strings.ReplaceAll(msg, "%v", "")
		l.gLog.Info(msg, zap.Any("data", data), zap.String("line", utils.FileWithLineNum()), requestIDField(ctx, l.requestIDKey))
	}
}

// Warn print warn messages
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		msg = strings.ReplaceAll(msg, "%v", "")
		l.gLog.Warn(msg, zap.Any("data", data), zap.String("line", utils.FileWithLineNum()), requestIDField(ctx, l.requestIDKey))
	}
}

// Error print error messages
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		msg = strings.ReplaceAll(msg, "%v", "")
		l.gLog.Warn(msg, zap.Any("data", data), zap.String("line", utils.FileWithLineNum()), requestIDField(ctx, l.requestIDKey))
	}
}

// Trace print sql message
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	var rowsField zap.Field
	if rows == -1 {
		rowsField = zap.String("rows", "-")
	} else {
		rowsField = zap.Int64("rows", rows)
	}

	var fileLineField zap.Field
	fileLine := utils.FileWithLineNum()
	ss := strings.Split(fileLine, "/internal/")
	if len(ss) == 2 {
		fileLineField = zap.String("file_line", ss[1])
	} else {
		fileLineField = zap.String("file_line", fileLine)
	}

	if err != nil {
		l.gLog.Warn("Gorm msg",
			zap.Error(err),
			zap.String("sql", sql),
			rowsField,
			zap.Float64("ms", float64(elapsed.Nanoseconds())/1e6),
			fileLineField,
			requestIDField(ctx, l.requestIDKey),
		)
		return
	}

	if l.logLevel >= logger.Info {
		l.gLog.Info("Gorm msg",
			zap.String("sql", sql),
			rowsField,
			zap.Float64("ms", float64(elapsed.Nanoseconds())/1e6),
			fileLineField,
			requestIDField(ctx, l.requestIDKey),
		)
		return
	}

	if l.logLevel >= logger.Warn {
		l.gLog.Warn("Gorm msg",
			zap.String("sql", sql),
			rowsField,
			zap.Float64("ms", float64(elapsed.Nanoseconds())/1e6),
			fileLineField,
			requestIDField(ctx, l.requestIDKey),
		)
	}
}

func requestIDField(ctx context.Context, requestIDKey string) zap.Field {
	if requestIDKey == "" {
		return zap.Skip()
	}

	var field zap.Field
	if requestIDKey != "" {
		if v, ok := ctx.Value(requestIDKey).(string); ok {
			field = zap.String(requestIDKey, v)
		} else {
			field = zap.Skip()
		}
	}
	return field
}
