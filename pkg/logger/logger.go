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
	if !isSave { // 在终端打印
		defaultLogger, err = log2Terminal(levelName, encoding)
		Infof("initialize logger finish, config is output to 'terminal', format=%s, level=%s", encoding, levelName)
	} else {
		defaultLogger = log2File(encoding, levelName, o.fileConfig)
		Info("initialize logger finish, config is output to 'file'", String("format", encoding), String("level", levelName), Any("file", o.fileConfig))
	}

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

// InitLogger 初始化日志
//
//	isSave 是否输出到文件，true: 是，false:输出到控台
//	filename 保存日志路径，例如："out.log"
//	level  输出日志级别 DEBUG, INFO, WARN, ERROR
//	encoding  输出格式 json:显示数据格式为json，console:显示数据格式为console(默认)
//		以console数据格式输出到控台，eg: InitLogger(false, "", "debug")
//		以json数据格式输出到控台，eg: InitLogger(false, "", "debug", "json")
//		以json数据格式输出到文件，eg: InitLogger(true, "out.log", "debug")
//
// Deprecated: use Init() instead.
func InitLogger(isSave bool, filename string, level string, encodingType ...string) error {
	// 保存日志路径
	if isSave && filename == "" {
		filename = "out.log" // 默认
	}

	// 日志输出等级
	levelName := strings.ToUpper(level)
	switch levelName {
	case levelDebug, levelInfo, levelWarn, levelError:
	default:
		fmt.Printf("unknown levelName: %s, use default: %s\n", levelName, levelDebug)
		levelName = levelDebug // 默认
	}

	var encoding string
	var js string
	if isSave { // 日志保存到文件
		encoding = formatJSON // 当日志输出到文件时，只有json格式
		js = fmt.Sprintf(`{
      		"level": "%s",
      		"encoding": "%s",
      		"outputPaths": ["%s"],
      		"errorOutputPaths": ["%s"]
      	}`, levelName, encoding, filename, filename)
	} else { // 在控台输出日志
		if len(encodingType) > 0 && encodingType[0] == formatJSON { // 控台模式下可以输出json格式，也可以输出console模式
			encoding = formatJSON
		} else {
			encoding = formatConsole
		}

		js = fmt.Sprintf(`{
      		"level": "%s",
            "encoding": "%s",
      		"outputPaths": ["stdout"],
            "errorOutputPaths": ["stdout"]
		}`, levelName, encoding)
	}

	var config zap.Config
	err := json.Unmarshal([]byte(js), &config)
	if err != nil {
		return err
	}

	config.EncoderConfig = zap.NewProductionEncoderConfig()

	config.EncoderConfig.EncodeTime = timeFormatter // 默认时间格式
	if isSave {
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	defaultLogger, err = config.Build()
	if err != nil {
		return err
	}

	// 打印log配置结果
	if isSave {
		Infof("initialize logger finish, base config is isSave=%t, filename=%s, level=%s, encoding=%s", isSave, filename, level, encoding)
	} else {
		Infof("initialize logger finish, base config is isSave=%t, level=%s, encoding=%s", isSave, level, encoding)
	}

	return nil
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
		//Warn("not yet initialized the log, use default log")
	}
}
