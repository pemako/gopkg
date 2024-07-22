package ctxlog

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var (
	svcLogger *CtxLogger
	loggers   map[string]*CtxLogger
)

// NewLogger sugar logger支持如 sprintf 输出的Infof，非结构化输入的 Infow, 性能上不如zap.Logger, 但是更方便
func NewLogger(conf Config) *CtxLogger {
	return NewLoggerWithExposedKey(conf, []string{})
}

func NewLoggerWithExposedKey(cfg Config, exposedKey []string) *CtxLogger {
	l := zap.New(newRotate(&cfg, false),
		zap.AddCaller(),
		zap.Fields(zap.String("svr", cfg.ServiceName)),
		zap.AddCallerSkip(2))

	tl := zap.New(newRotate(&cfg, true),
		zap.AddCaller(),
		zap.Fields(zap.String("svr", cfg.ServiceName)),
		zap.AddCallerSkip(2))

	if l != nil && tl != nil {
		sugar := l.Sugar()
		sugarT := tl.Sugar()
		if sugar != nil && sugarT != nil {
			return NewCtxLogger(sugar, sugarT, exposedKey)
		}
	}

	return nil
}

// NewStdoutLogger 创建一个输出到stdout的日志logger
// @serviceName 服务名称，应当与当前模块名称保持一致
// @loggerName 这个logger的名称，如果为空则默认为“main”,表明为主日志, 因为stdout模式所有的日志都会输出到stdout中，loggerName有助于区分日志logger
// @formatType json/custom 格式的日志
// @level 日志级别，参考config.go
// @ctxFunc 见bfo_ctx_logger type CtxFunc func(ctx context.Context) map[string]any, 一个解析ctx中数据的函数,用于日志记录
func NewStdoutLogger(serviceName string, loggerName string, formatType string, level string) *CtxLogger {
	return NewStdoutLoggerWithExposedKey(serviceName, loggerName, formatType, level, []string{})
}

func NewStdoutLoggerWithExposedKey(serviceName string, loggerName string, formatType string, level string, exposedKey []string) *CtxLogger {
	if loggerName == "" {
		loggerName = "main"
	}

	sugar := zap.New(
		newStdout(formatType, level),
		zap.AddCaller(),
		zap.Development(),
		zap.Fields(zap.String("svc", serviceName), zap.String("loggerName", loggerName)),
		zap.AddCallerSkip(1),
	).Sugar()

	return NewCtxLogger(sugar, sugar, exposedKey)
}

func SetLogger(logger *CtxLogger) {
	svcLogger = logger
	// 程序退出时，通过调用logger.Sync()输出buffer中的内容
	go func() {
		c := make(chan os.Signal, 5)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		for {
			sig := <-c
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				signal.Stop(c)
				svcLogger.Sync()
			default:
				continue
			}
		}
	}()
}

func GetLogger() *CtxLogger {
	return svcLogger
}

func SetLoggerByName(name string, logger *CtxLogger) {
	loggers[name] = logger
	// 程序退出时，通过调用logger.Sync()输出buffer中的内容
	go func() {
		c := make(chan os.Signal, 5)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		for {
			sig := <-c
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				signal.Stop(c)
				logger.Sync()
			default:
				continue
			}
		}
	}()
}

func GetLoggerByName(name string) *CtxLogger {
	if logger, ok := loggers[name]; ok {
		return logger
	}

	return nil
}

func checkCtxServiceLoggerNotNull() bool {
	if svcLogger == nil {
		fmt.Fprint(os.Stderr, "ctx service logger not init! log will not be printed!!\n")
		return false
	}
	return true
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(ctx context.Context, msg string) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.debug(ctx, msg)
	}
}

// Info uses fmt.Sprint to construct and log a message.
func Info(ctx context.Context, msg string) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.info(ctx, msg)
	}
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(ctx context.Context, msg string) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.warn(ctx, msg)
	}
}

// Error uses fmt.Sprint to construct and log a message.
func Error(ctx context.Context, msg string) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.error(ctx, msg)
	}
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(ctx context.Context, msg string) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.dPanic(ctx, msg)
	}
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(ctx context.Context, msg string) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.panic(ctx, msg)
	}
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(ctx context.Context, msg string) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.fatal(ctx, msg)
	}
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(ctx context.Context, template string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.debugf(ctx, template, args...)
	}
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(ctx context.Context, template string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.infof(ctx, template, args...)
	}
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(ctx context.Context, template string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.warnf(ctx, template, args...)
	}
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(ctx context.Context, template string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.errorf(ctx, template, args...)
	}
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(ctx context.Context, template string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.dPanicf(ctx, template, args...)
	}
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(ctx context.Context, template string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.panicf(ctx, template, args...)
	}
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(ctx context.Context, template string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.fatalf(ctx, template, args...)
	}
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func Debugw(ctx context.Context, msg string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.debugw(ctx, msg, args...)
	}
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(ctx context.Context, msg string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.infow(ctx, msg, args...)
	}

}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(ctx context.Context, msg string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.warnw(ctx, msg, args...)
	}
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(ctx context.Context, msg string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.errorw(ctx, msg, args...)
	}
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(ctx context.Context, msg string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.dPanicw(ctx, msg, args...)
	}
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(ctx context.Context, msg string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.panicw(ctx, msg, args...)
	}
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(ctx context.Context, msg string, args ...any) {
	if checkCtxServiceLoggerNotNull() {
		svcLogger.fatalw(ctx, msg, args...)
	}
}

func init() {
	loggers = make(map[string]*CtxLogger)
}
