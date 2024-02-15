package gocron

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var defaultLog, _ = zap.NewProduction()

type options struct {
	zapLog *zap.Logger
}

func defaultOptions() *options {
	return &options{
		zapLog: defaultLog,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option set the cron options.
type Option func(*options)

// WithLog set log
func WithLog(log *zap.Logger) Option {
	return func(o *options) {
		o.zapLog = log
	}
}

type zapLog struct {
	zapLog *zap.Logger
}

// Info print info
func (l *zapLog) Info(msg string, keysAndValues ...interface{}) {
	if msg == "wake" { // 忽略wake
		return
	}
	msg = "cron_" + msg
	fields := parseKVs(keysAndValues)
	l.zapLog.Info(msg, fields...)
}

// Error print error
func (l *zapLog) Error(err error, msg string, keysAndValues ...interface{}) {
	fields := parseKVs(keysAndValues)
	fields = append(fields, zap.String("err", err.Error()))
	msg = "cron_" + msg
	l.zapLog.Error(msg, fields...)
}

func parseKVs(kvs interface{}) []zap.Field {
	var fields []zap.Field

	infos, ok := kvs.([]interface{})
	if !ok {
		return fields
	}

	l := len(infos)
	if l%2 == 1 {
		return fields
	}

	for i := 0; i < l; i += 2 {
		key := infos[i].(string) //nolint
		value := infos[i+1]

		// replace id with task name
		if key == "entry" {
			if id, ok := value.(cron.EntryID); ok {
				key = "task"
				if v, isExist := idName.Load(id); isExist {
					value = v
				}
			}
		}

		fields = append(fields, zap.Any(key, value))
	}

	return fields
}
