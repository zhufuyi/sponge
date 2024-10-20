package gocron

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	SecondType = 0
	MinuteType = 1
)

var defaultLog, _ = zap.NewProduction()

type options struct {
	zapLog           *zap.Logger
	isOnlyPrintError bool // default false

	granularity int // 0: second, 1: minute
}

func defaultOptions() *options {
	return &options{
		zapLog:           defaultLog,
		isOnlyPrintError: false,

		granularity: SecondType,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option set the cron options.
type Option func(*options)

// WithGranularity set log
func WithGranularity(granularity int) Option {
	return func(o *options) {
		if granularity >= MinuteType {
			granularity = MinuteType
		} else {
			granularity = SecondType
		}
		o.granularity = granularity
	}
}

// WithLog set granularity
func WithLog(log *zap.Logger, isOnlyPrintError ...bool) Option {
	return func(o *options) {
		if len(isOnlyPrintError) > 0 {
			o.isOnlyPrintError = isOnlyPrintError[0]
		}
		o.zapLog = log
	}
}

type zapLog struct {
	zapLog           *zap.Logger
	isOnlyPrintError bool
}

// Info print info
func (l *zapLog) Info(msg string, keysAndValues ...interface{}) {
	if l.zapLog == nil || l.isOnlyPrintError {
		return
	}
	if msg == "wake" { // 忽略wake
		return
	}
	msg = "cron_" + msg
	fields := parseKVs(keysAndValues)
	l.zapLog.Info(msg, fields...)
}

// Error print error
func (l *zapLog) Error(err error, msg string, keysAndValues ...interface{}) {
	if l.zapLog == nil {
		return
	}
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
