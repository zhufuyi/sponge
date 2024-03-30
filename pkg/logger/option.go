package logger

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

var (
	defaultLevel    = "debug" // output log levels debug, info, warn, error, default is debug
	defaultEncoding = formatConsole
	defaultIsSave   = false // false:output to terminal, true:output to file, default is false

	defaultFilename      = "out.log" // file name
	defaultMaxSize       = 10        // maximum file size (MB)
	defaultMaxBackups    = 100       // maximum number of old files
	defaultMaxAge        = 30        // maximum number of days for old documents
	defaultIsCompression = false     // whether to compress and archive old files
)

type options struct {
	level    string
	encoding string
	isSave   bool

	fileConfig *fileOptions

	hooks []func(zapcore.Entry) error
}

func defaultOptions() *options {
	return &options{
		level:    defaultLevel,
		encoding: defaultEncoding,
		isSave:   defaultIsSave,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option set the logger options.
type Option func(*options)

// WithLevel setting the log level
func WithLevel(levelName string) Option {
	return func(o *options) {
		levelName = strings.ToUpper(levelName)
		switch levelName {
		case levelDebug, levelInfo, levelWarn, levelError:
			o.level = levelName
		default:
			o.level = levelDebug
		}
	}
}

// WithFormat set the output log format, console or json
func WithFormat(format string) Option {
	return func(o *options) {
		if strings.ToLower(format) == formatJSON {
			o.encoding = formatJSON
		}
	}
}

// WithSave save log to file
func WithSave(isSave bool, opts ...FileOption) Option {
	return func(o *options) {
		if isSave {
			o.isSave = true
			fo := defaultFileOptions()
			fo.apply(opts...)
			o.fileConfig = fo
		}
	}
}

// WithHooks set the log hooks
func WithHooks(hooks ...func(zapcore.Entry) error) Option {
	return func(o *options) {
		o.hooks = hooks
	}
}

// ------------------------------------------------------------------------------------------

type fileOptions struct {
	filename      string
	maxSize       int
	maxBackups    int
	maxAge        int
	isCompression bool
}

func defaultFileOptions() *fileOptions {
	return &fileOptions{
		filename:      defaultFilename,
		maxSize:       defaultMaxSize,
		maxBackups:    defaultMaxBackups,
		maxAge:        defaultMaxAge,
		isCompression: defaultIsCompression,
	}
}

func (o *fileOptions) apply(opts ...FileOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// FileOption set the file options.
type FileOption func(*fileOptions)

// WithFileName set log filename
func WithFileName(filename string) FileOption {
	return func(f *fileOptions) {
		if filename != "" {
			f.filename = filename
		}
	}
}

// WithFileMaxSize set maximum file size (MB)
func WithFileMaxSize(maxSize int) FileOption {
	return func(f *fileOptions) {
		if f.maxSize > 0 {
			f.maxSize = maxSize
		}
	}
}

// WithFileMaxBackups set maximum number of old files
func WithFileMaxBackups(maxBackups int) FileOption {
	return func(f *fileOptions) {
		if f.maxBackups > 0 {
			f.maxBackups = maxBackups
		}
	}
}

// WithFileMaxAge set maximum number of days for old documents
func WithFileMaxAge(maxAge int) FileOption {
	return func(f *fileOptions) {
		if f.maxAge > 0 {
			f.maxAge = maxAge
		}
	}
}

// WithFileIsCompression set whether to compress log files
func WithFileIsCompression(isCompression bool) FileOption {
	return func(f *fileOptions) {
		f.isCompression = isCompression
	}
}
