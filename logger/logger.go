package logger

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pemako/gopkg/lumberjack"
	"github.com/pemako/gopkg/rotatelogs"
)

var (
	svcLogger *zap.SugaredLogger
)

func New(cfg Config) *zap.SugaredLogger {
	logger := zap.New(newRotate(&cfg), zap.AddCaller(), zap.AddCallerSkip(1))
	if cfg.EnableStackTrace {
		logger = logger.WithOptions(zap.AddStacktrace(strToZapLevel(cfg.Level)))
	}

	lg := logger.Sugar()

	setLogger(lg)
	return lg
}

func GetLogger() *zap.SugaredLogger {
	return svcLogger
}

func setLogger(logger *zap.SugaredLogger) {
	svcLogger = logger
	go func() {
		c := make(chan os.Signal, 5)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		for {
			sig := <-c
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				signal.Stop(c)
				_ = svcLogger.Sync()
			default:
				continue
			}
		}
	}()
}

func newRotate(cfg *Config) zapcore.Core {
	encoderCfg := encoderConfig()
	logLevels := genLogLevels(cfg.FileName)
	var encoder zapcore.Encoder
	switch cfg.FormatType {
	case JsonFormat:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	var cores []zapcore.Core
	_lv := strToZapLevel(cfg.Level)
	for lv, lg := range logLevels {
		if _lv <= lv {
			fileName := cfg.FilePath + "/" + lg
			var core zapcore.Core
			switch cfg.RotateType {
			case "size":
				hook := sizeRotate(fileName, cfg.SizeRotate)
				core = zapcore.NewCore(encoder, zapcore.AddSync(&hook), genLevelEnabler(lv))
			default:
				hook := timeRotate(fileName, cfg.TimeRotate)
				core = zapcore.NewCore(encoder, zapcore.AddSync(hook), genLevelEnabler(lv))
			}
			cores = append(cores, core)
		}
	}

	return zapcore.NewTee(cores...)
}

func sizeRotate(name string, cfg *SizeRotate) lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   name,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  false,
		Compress:   false,
	}
}

func timeRotate(name string, cfg *TimeRotate) io.Writer {
	hook, err := rotatelogs.New(
		name+"."+cfg.Format,
		rotatelogs.WithLinkName(name),
		rotatelogs.WithMaxAge(time.Hour*time.Duration(cfg.MaxAge)),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(cfg.RotateTime)),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func genLogLevels(name string) map[zapcore.Level]string {
	return map[zapcore.Level]string{
		zapcore.DebugLevel: name + ".debug.log",
		zapcore.InfoLevel:  name + ".info.log",
		zapcore.WarnLevel:  name + ".warn.log",
		zapcore.ErrorLevel: name + ".error.log",
	}
}

func genLevelEnabler(lv zapcore.Level) zapcore.LevelEnabler {
	switch lv {
	case zapcore.DebugLevel:
		return zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.DebugLevel
		})
	case zapcore.InfoLevel:
		return zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.InfoLevel
		})
	case zapcore.WarnLevel:
		return zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.WarnLevel
		})
	case zapcore.ErrorLevel:
		return zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		})
	default:
		atomicLevel := zap.NewAtomicLevel()
		atomicLevel.SetLevel(lv)
		return atomicLevel
	}
}
