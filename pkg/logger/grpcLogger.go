package logger

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

type grpcLogger struct {
	zLog      *zap.Logger
	verbosity int
}

// ReplaceGRPCLoggerV2 replace grpc logger v2
func ReplaceGRPCLoggerV2(l *zap.Logger) {
	zzl := &grpcLogger{
		zLog:      l.With(zap.String("log_from", "grpc_system")),
		verbosity: 0,
	}
	grpclog.SetLoggerV2(zzl)
}

func (l *grpcLogger) Info(args ...interface{}) {
	l.zLog.Info(fmt.Sprint(args...))
}

func (l *grpcLogger) Infoln(args ...interface{}) {
	l.zLog.Info(fmt.Sprint(args...))
}

func (l *grpcLogger) Infof(format string, args ...interface{}) {
	l.zLog.Info(fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Warning(args ...interface{}) {
	l.zLog.Warn(fmt.Sprint(args...))
}

func (l *grpcLogger) Warningln(args ...interface{}) {
	l.zLog.Warn(fmt.Sprint(args...))
}

func (l *grpcLogger) Warningf(format string, args ...interface{}) {
	l.zLog.Warn(fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Error(args ...interface{}) {
	l.zLog.Error(fmt.Sprint(args...))
}

func (l *grpcLogger) Errorln(args ...interface{}) {
	l.zLog.Error(fmt.Sprint(args...))
}

func (l *grpcLogger) Errorf(format string, args ...interface{}) {
	l.zLog.Error(fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Fatal(args ...interface{}) {
	l.zLog.Fatal(fmt.Sprint(args...))
}

func (l *grpcLogger) Fatalln(args ...interface{}) {
	l.zLog.Fatal(fmt.Sprint(args...))
}

func (l *grpcLogger) Fatalf(format string, args ...interface{}) {
	l.zLog.Fatal(fmt.Sprintf(format, args...))
}

func (l *grpcLogger) V(level int) bool {
	return l.verbosity <= level
}
