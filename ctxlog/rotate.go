package ctxlog

import (
	"io"
	"os"
	"time"

	"github.com/pemako/gopkg/lumberjack"
	"github.com/pemako/gopkg/rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newRotate(cfg *Config, isTrack bool) zapcore.Core {
	encoderCfg := encoderConfig()
	if cfg.FileName == "" {
		cfg.FileName = cfg.ServiceName
	}
	logLevels := genLogLevels(cfg.FileName)
	if isTrack {
		encoderCfg = encoderConfigForTrack()
		logLevels = genLogLevelsForTrack(cfg.FileName)
	}

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
			case RotateTypeSize:
				hook := sizeRotate(fileName, cfg.SizeRotateConfig)
				core = zapcore.NewCore(encoder, zapcore.AddSync(&hook), genLevelEnabler(lv))
			default:
				hook := timeRotateHook(fileName, cfg.TimeRotateConfig)
				core = zapcore.NewCore(encoder, zapcore.AddSync(hook), genLevelEnabler(lv))
			}
			cores = append(cores, core)
		}
	}

	return zapcore.NewTee(cores...)
}

func newStdout(formatType string, level string) zapcore.Core {
	ec := encoderConfig()
	var encoder zapcore.Encoder
	if formatType == JsonFormat {
		encoder = zapcore.NewJSONEncoder(ec)
	} else {
		encoder = zapcore.NewConsoleEncoder(ec)
	}
	return zapcore.NewCore(encoder, os.Stdout, strToZapLevel(level))
}

func sizeRotate(name string, cfg *SizeRotateConfig) lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   name,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  false,
		Compress:   cfg.Compress,
	}
}

func timeRotateHook(name string, cfg *TimeRotateConfig) io.Writer {
	hook, err := rotatelogs.New(
		name+"."+cfg.FileNameFormat,
		rotatelogs.WithLinkName(name),
		rotatelogs.WithMaxAge(time.Hour*time.Duration(cfg.MaxAge)),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(cfg.RotateTime)),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func genLogLevelsForTrack(name string) map[zapcore.Level]string {
	return map[zapcore.Level]string{
		zapcore.InfoLevel: name + ".track.log",
	}
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
