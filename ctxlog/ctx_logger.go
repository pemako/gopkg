package ctxlog

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type CtxFunc func(ctx context.Context) map[string]any

type CtxLogger struct {
	logger      *zap.SugaredLogger
	trackLogger *zap.SugaredLogger
	exposedKey  []string
}

func NewCtxLogger(l, t *zap.SugaredLogger, k []string) *CtxLogger {
	return &CtxLogger{
		logger:      l,
		trackLogger: t,
		exposedKey:  k,
	}
}

func checkExistInList(k string, kl []string) bool {
	for _, v := range kl {
		if k == v {
			return true
		}
	}
	return false
}

func (l *CtxLogger) appendCtxArgs(ctx context.Context, args ...any) []any {
	ext := []any{}
	data := Ctx{StartTime: time.Now()}
	if ctx != nil {
		if d, exist := GetCtxData(ctx); exist {
			data = d
		}
	}

	cost := time.Since(data.StartTime).Milliseconds()
	ext = append(ext, "uuid", data.UUID)
	ext = append(ext, "infc", data.Infc)
	ext = append(ext, "req", data.Req)
	ext = append(ext, "cost", cost)

	for k, v := range data.Ext {
		if checkExistInList(k, l.exposedKey) {
			ext = append(ext, k, v)
		}
	}

	args = append(args, ext...)
	return args
}

func (l *CtxLogger) Debug(ctx context.Context, msg string) {
	l.debug(ctx, msg)
}

func (l *CtxLogger) Info(ctx context.Context, msg string) {
	l.info(ctx, msg)
}

func (l *CtxLogger) Warn(ctx context.Context, msg string) {
	l.warn(ctx, msg)
}

func (l *CtxLogger) Error(ctx context.Context, msg string) {
	l.error(ctx, msg)
}

func (l *CtxLogger) DPanic(ctx context.Context, msg string) {
	l.dPanic(ctx, msg)
}

func (l *CtxLogger) Panic(ctx context.Context, msg string) {
	l.panic(ctx, msg)
}

func (l *CtxLogger) Fatal(ctx context.Context, msg string) {
	l.fatal(ctx, msg)
}

func (l *CtxLogger) Debugf(ctx context.Context, template string, args ...any) {
	l.debugf(ctx, template, args...)
}

func (l *CtxLogger) Infof(ctx context.Context, template string, args ...any) {
	l.infof(ctx, template, args...)
}

func (l *CtxLogger) Warnf(ctx context.Context, template string, args ...any) {
	l.warnf(ctx, template, args...)
}

func (l *CtxLogger) Errorf(ctx context.Context, template string, args ...any) {
	l.errorf(ctx, template, args...)
}

func (l *CtxLogger) DPanicf(ctx context.Context, template string, args ...any) {
	l.dPanicf(ctx, template, args...)
}

func (l *CtxLogger) Panicf(ctx context.Context, template string, args ...any) {
	l.panicf(ctx, template, args...)
}

func (l *CtxLogger) Fatalf(ctx context.Context, template string, args ...any) {
	l.fatalf(ctx, template, args...)
}

func (l *CtxLogger) Debugw(ctx context.Context, msg string, args ...any) {
	l.debugw(ctx, msg, args...)
}

func (l *CtxLogger) Infow(ctx context.Context, msg string, args ...any) {
	l.infow(ctx, msg, args...)
}

func (l *CtxLogger) Warnw(ctx context.Context, msg string, args ...any) {
	l.warnw(ctx, msg, args...)
}

func (l *CtxLogger) Errorw(ctx context.Context, msg string, args ...any) {
	l.errorw(ctx, msg, args...)
}

func (l *CtxLogger) DPanicw(ctx context.Context, msg string, args ...any) {
	l.dPanicw(ctx, msg, args...)
}

func (l *CtxLogger) Panicw(ctx context.Context, msg string, args ...any) {
	l.panicw(ctx, msg, args...)
}

func (l *CtxLogger) Fatalw(ctx context.Context, msg string, args ...any) {
	l.fatalw(ctx, msg, args...)
}

func (l *CtxLogger) debug(ctx context.Context, msg string) {
	args := l.appendCtxArgs(ctx)
	l.logger.Debugw(msg, args...)
}

// Info uses fmt.Sprint to construct and log a message.
func (l *CtxLogger) info(ctx context.Context, msg string) {
	args := l.appendCtxArgs(ctx)
	l.logger.Infow(msg, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *CtxLogger) warn(ctx context.Context, msg string) {
	args := l.appendCtxArgs(ctx)
	l.logger.Warnw(msg, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *CtxLogger) error(ctx context.Context, msg string) {
	args := l.appendCtxArgs(ctx)
	l.logger.Errorw(msg, args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *CtxLogger) dPanic(ctx context.Context, msg string) {
	args := l.appendCtxArgs(ctx)
	l.logger.DPanicw(msg, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *CtxLogger) panic(ctx context.Context, msg string) {
	args := l.appendCtxArgs(ctx)
	l.logger.Panicw(msg, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *CtxLogger) fatal(ctx context.Context, msg string) {
	args := l.appendCtxArgs(ctx)
	l.logger.Fatalw(msg, args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *CtxLogger) debugf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	args = l.appendCtxArgs(ctx)
	l.logger.Debugw(msg, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *CtxLogger) infof(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	args = l.appendCtxArgs(ctx)
	l.logger.Infow(msg, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *CtxLogger) warnf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	args = l.appendCtxArgs(ctx)
	l.logger.Warnw(msg, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *CtxLogger) errorf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	args = l.appendCtxArgs(ctx)
	l.logger.Errorw(msg, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *CtxLogger) dPanicf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	args = l.appendCtxArgs(ctx)
	l.logger.DPanicw(msg, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *CtxLogger) panicf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	args = l.appendCtxArgs(ctx)
	l.logger.Panicw(msg, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *CtxLogger) fatalf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	args = l.appendCtxArgs(ctx)
	l.logger.Fatalw(msg, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func (l *CtxLogger) debugw(ctx context.Context, msg string, args ...any) {
	ctxArgs := l.appendCtxArgs(ctx)
	args = append(ctxArgs, args...)
	l.logger.Debugw(msg, args...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *CtxLogger) infow(ctx context.Context, msg string, args ...any) {
	ctxArgs := l.appendCtxArgs(ctx)
	args = append(ctxArgs, args...)
	l.logger.Infow(msg, args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *CtxLogger) warnw(ctx context.Context, msg string, args ...any) {
	ctxArgs := l.appendCtxArgs(ctx)
	args = append(ctxArgs, args...)
	l.logger.Warnw(msg, args...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *CtxLogger) errorw(ctx context.Context, msg string, args ...any) {
	ctxArgs := l.appendCtxArgs(ctx)
	args = append(ctxArgs, args...)
	l.logger.Errorw(msg, args...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (l *CtxLogger) dPanicw(ctx context.Context, msg string, args ...any) {
	ctxArgs := l.appendCtxArgs(ctx)
	args = append(ctxArgs, args...)
	l.logger.DPanicw(msg, args...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (l *CtxLogger) panicw(ctx context.Context, msg string, args ...any) {
	ctxArgs := l.appendCtxArgs(ctx)
	args = append(ctxArgs, args...)
	l.logger.Panicw(msg, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (l *CtxLogger) fatalw(ctx context.Context, msg string, args ...any) {
	ctxArgs := l.appendCtxArgs(ctx)
	args = append(ctxArgs, args...)
	l.logger.Fatalw(msg, args...)
}

// Sync flushes any buffered log entries.
func (l *CtxLogger) Sync() error {
	return l.logger.Sync()
}
