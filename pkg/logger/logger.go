package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	formatConsole = "console"
	formatJSON    = "json"

	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelWarn  = "WARN"
	levelError = "ERROR"
)

var defaultLogger *zap.Logger

func getLogger() *zap.Logger {
	checkNil()
	return defaultLogger.WithOptions(zap.AddCallerSkip(1))
}

// Init 初始化打印设置
// 在终端打印debug级别日志示例：Init()
// 在终端打印info级别日志示例：Init(WithLevel("info"))
// 在终端打印json格式、debug级别日志示例：Init(WithFormat("json"))
// 把日志输出到文件out.log，使用默认的切割日志相关参数，debug级别日志示例：Init(WithSave())
// 把日志输出到指定文件，自定义设置日志文件切割日志参数，json格式，debug级别日志示例：
// Init(
//
//	  WithFormat("json"),
//	  WithSave(true,
//
//			WithFileName("my.log"),
//			WithFileMaxSize(5),
//			WithFileMaxBackups(5),
//			WithFileMaxAge(10),
//			WithFileIsCompression(true),
//		))
func Init(opts ...Option) (*zap.Logger, error) {
	o := defaultOptions()
	o.apply(opts...)
	isSave := o.isSave
	levelName := o.level
	encoding := o.encoding

	var err error
	var zapLog *zap.Logger
	var str string
	if !isSave { // 在终端打印
		zapLog, err = log2Terminal(levelName, encoding)
		if err != nil {
			panic(err)
		}
		str = fmt.Sprintf("initialize logger finish, config is output to 'terminal', format=%s, level=%s", encoding, levelName)
	} else {
		zapLog = log2File(encoding, levelName, o.fileConfig)
		str = fmt.Sprintf("initialize logger finish, config is output to 'file', format=%s, level=%s, file=%s", encoding, levelName, o.fileConfig.filename)
	}

	defaultLogger = zapLog
	Info(str)

	return defaultLogger, err
}

// 在终端打印
func log2Terminal(levelName string, encoding string) (*zap.Logger, error) {
	js := fmt.Sprintf(`{
      		"level": "%s",
            "encoding": "%s",
      		"outputPaths": ["stdout"],
            "errorOutputPaths": ["stdout"]
		}`, levelName, encoding)

	var config zap.Config
	err := json.Unmarshal([]byte(js), &config)
	if err != nil {
		return nil, err
	}

	config.EncoderConfig = zap.NewProductionEncoderConfig()
	config.EncoderConfig.EncodeTime = timeFormatter                // 默认时间格式
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 在日志文件中使用大写字母记录日志级别
	return config.Build()
}

// 输出到文件
func log2File(encoding string, levelName string, fo *fileOptions) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 修改时间编码器
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 在日志文件中使用大写字母记录日志级别
	var encoder zapcore.Encoder
	if encoding == formatConsole { // console格式
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else { // json格式
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fo.filename,      // 日志文件的位置；
		MaxSize:    fo.maxSize,       // 在进行切割之前，日志文件的最大值(单位MB)
		MaxBackups: fo.maxBackups,    // 保留旧文件的最大个数
		MaxAge:     fo.maxAge,        // 保留旧文件的最大天数
		Compress:   fo.isCompression, // 是否压缩/归档旧文件
	})
	core := zapcore.NewCore(encoder, ws, getLevelSize(levelName))

	// zap.AddCaller()  添加将调用函数信息记录到日志中的功能。
	return zap.New(core, zap.AddCaller())
}

// DEBUG(默认), INFO, WARN, ERROR
func getLevelSize(levelName string) zapcore.Level {
	levelName = strings.ToUpper(levelName)
	switch levelName {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarn:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	}
	return zapcore.DebugLevel
}

func timeFormatter(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// GetWithSkip 获取defaultLogger，设置跳过的caller值，自定义显示代码行数
func GetWithSkip(skip int) *zap.Logger {
	checkNil()
	return defaultLogger.WithOptions(zap.AddCallerSkip(skip))
}

// Get 获取logger对象
func Get() *zap.Logger {
	checkNil()
	return defaultLogger
}

func checkNil() {
	if defaultLogger == nil {
		_, err := Init() // 默认输出到控台
		if err != nil {
			panic(err)
		}
	}
}
