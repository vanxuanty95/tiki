package logger

import (
	"context"
	"log"
)

const (
	INFO_LEVEL  = "info"
	DEBUG_LEVEL = "debug"
	WARN_LEVEL  = "warn"

	TIKILoggerConText = "tikiLogger"
)

type (
	Logger interface {
		WithFields(fields map[string]interface{}) Logger
		WithPrefix(prefix string) Logger

		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Printf(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Panicf(format string, args ...interface{})

		Debug(args ...interface{})
		Info(args ...interface{})
		Print(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Panic(args ...interface{})

		Debugln(args ...interface{})
		Infoln(args ...interface{})
		Println(args ...interface{})
		Warnln(args ...interface{})
		Errorln(args ...interface{})
		Panicln(args ...interface{})

		TIKIDebugf(ctx context.Context, format string, args ...interface{})
		TIKIInfof(ctx context.Context, format string, args ...interface{})
		TIKIPrintf(ctx context.Context, format string, args ...interface{})
		TIKIWarnf(ctx context.Context, format string, args ...interface{})
		TIKIErrorf(ctx context.Context, format string, args ...interface{})
		TIKIPanicf(ctx context.Context, format string, args ...interface{})

		TIKIDebug(ctx context.Context, args ...interface{})
		TIKIInfo(ctx context.Context, args ...interface{})
		TIKIPrint(ctx context.Context, args ...interface{})
		TIKIWarn(ctx context.Context, args ...interface{})
		TIKIError(ctx context.Context, args ...interface{})
		TIKIPanic(ctx context.Context, args ...interface{})

		TIKIDebugln(ctx context.Context, args ...interface{})
		TIKIInfoln(ctx context.Context, args ...interface{})
		TIKIPrintln(ctx context.Context, args ...interface{})
		TIKIWarnln(ctx context.Context, args ...interface{})
		TIKIErrorln(ctx context.Context, args ...interface{})
		TIKIPanicln(ctx context.Context, args ...interface{})
	}
)

var std Logger

func init() {
	var err error = nil
	std, err = newLogrusLogger(INFO_LEVEL)
	if err != nil {
		log.Panic(err)
	}
}

func WithFields(fields map[string]interface{}) Logger {
	return std.WithFields(fields)
}

func WithPrefix(prefix string) Logger {
	return std.WithPrefix(prefix)
}

func WithMetricType(metricType string) Logger {
	return std.WithFields(map[string]interface{}{
		"metric_type": metricType,
	})
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	std.Debug(args...)
}
func Info(args ...interface{}) {
	std.Info(args...)
}
func Print(args ...interface{}) {
	std.Print(args...)
}
func Warn(args ...interface{}) {
	std.Warn(args...)
}
func Error(args ...interface{}) {
	std.Error(args...)
}
func Panic(args ...interface{}) {
	std.Panic(args...)
}

func Debugln(args ...interface{}) {
	std.Debugln(args...)
}
func Infoln(args ...interface{}) {
	std.Infoln(args...)
}
func Println(args ...interface{}) {
	std.Println(args...)
}
func Warnln(args ...interface{}) {
	std.Warnln(args...)
}
func Errorln(args ...interface{}) {
	std.Errorln(args...)
}
func Panicln(args ...interface{}) {
	std.Panicln(args...)
}

func TIKIDebugf(ctx context.Context, format string, args ...interface{}) {
	std.TIKIDebugf(ctx, format, args...)
}
func TIKIInfof(ctx context.Context, format string, args ...interface{}) {
	std.TIKIInfof(ctx, format, args...)
}
func TIKIPrintf(ctx context.Context, format string, args ...interface{}) {
	std.TIKIPrintf(ctx, format, args...)
}
func TIKIWarnf(ctx context.Context, format string, args ...interface{}) {
	std.TIKIWarnf(ctx, format, args...)
}
func TIKIErrorf(ctx context.Context, format string, args ...interface{}) {
	std.TIKIErrorf(ctx, format, args...)
}
func TIKIPanicf(ctx context.Context, format string, args ...interface{}) {
	std.TIKIPanicf(ctx, format, args...)
}

func TIKIDebug(ctx context.Context, args ...interface{}) {
	std.TIKIDebug(ctx, args...)
}
func TIKIInfo(ctx context.Context, args ...interface{}) {
	std.TIKIInfo(ctx, args...)
}
func TIKIPrint(ctx context.Context, args ...interface{}) {
	std.TIKIPrint(ctx, args...)
}
func TIKIWarn(ctx context.Context, args ...interface{}) {
	std.TIKIWarn(ctx, args...)
}
func TIKIError(ctx context.Context, args ...interface{}) {
	std.TIKIError(ctx, args...)
}
func TIKIPanic(ctx context.Context, args ...interface{}) {
	std.TIKIPanic(ctx, args...)
}

func TIKIDebugln(ctx context.Context, args ...interface{}) {
	std.TIKIDebugln(ctx, args...)
}
func TIKIInfoln(ctx context.Context, args ...interface{}) {
	std.TIKIInfoln(ctx, args...)
}
func TIKIPrintln(ctx context.Context, args ...interface{}) {
	std.TIKIPrintln(ctx, args...)
}
func TIKIWarnln(ctx context.Context, args ...interface{}) {
	std.TIKIWarnln(ctx, args...)
}
func TIKIErrorln(ctx context.Context, args ...interface{}) {
	std.TIKIErrorln(ctx, args...)
}
func TIKIPanicln(ctx context.Context, args ...interface{}) {
	std.TIKIPanicln(ctx, args...)
}
