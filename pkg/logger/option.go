package logger

import "strings"

var (
	defaultLevel    = "debug" //  输出日志级别 debug, info, warn, error，默认是debug
	defaultEncoding = formatConsole
	defaultIsSave   = false // false:输出到终端，true:输出到文件，默认是false

	// 保存文件相关默认设置
	defaultFilename      = "out.log" // 文件名称
	defaultMaxSize       = 10        // 最大文件大小(MB)
	defaultMaxBackups    = 100       // 保留旧文件的最大个数
	defaultMaxAge        = 30        // 保留旧文件的最大天数
	defaultIsCompression = false     // 是否压缩归档旧文件
)

type options struct {
	level    string
	encoding string
	isSave   bool

	// 保存文件相关默认设置
	fileConfig *fileOptions
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

// WithLevel 输出日志级别
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

// WithFormat 设置输出日志格式，console或json
func WithFormat(format string) Option {
	return func(o *options) {
		if strings.ToLower(format) == formatJSON {
			o.encoding = formatJSON
		}
	}
}

// WithSave 保存日志到指定文件
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

// WithFileName 自定义文件名称
func WithFileName(filename string) FileOption {
	return func(f *fileOptions) {
		if filename != "" {
			f.filename = filename
		}
	}
}

// WithFileMaxSize 自定义最大文件大小(MB)
func WithFileMaxSize(maxSize int) FileOption {
	return func(f *fileOptions) {
		if f.maxSize > 0 {
			f.maxSize = maxSize
		}
	}
}

// WithFileMaxBackups 自定义保留旧文件的最大个数
func WithFileMaxBackups(maxBackups int) FileOption {
	return func(f *fileOptions) {
		if f.maxBackups > 0 {
			f.maxBackups = maxBackups
		}
	}
}

// WithFileMaxAge 保留旧文件的最大天数
func WithFileMaxAge(maxAge int) FileOption {
	return func(f *fileOptions) {
		if f.maxAge > 0 {
			f.maxAge = maxAge
		}
	}
}

// WithFileIsCompression 自定义是否压缩归档旧文件
func WithFileIsCompression(isCompression bool) FileOption {
	return func(f *fileOptions) {
		f.isCompression = isCompression
	}
}
