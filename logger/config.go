package logger

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug = zapcore.DebugLevel
	LevelInfo  = zapcore.InfoLevel
	LevelWarn  = zapcore.WarnLevel
	LevelError = zapcore.ErrorLevel
)

const (
	JsonFormat    = "json"
	ConsoleFormat = "console"
)

type SizeRotate struct {
	MaxSize    int `json:"max_size" yaml:"maxSize"`
	MaxBackups int `json:"max_backups" yaml:"maxBackups"`
	MaxAge     int `json:"max_age" yaml:"maxAge"`
}

type TimeRotate struct {
	Format     string `json:"format" yaml:"format"`
	MaxAge     int    `json:"max_age" yaml:"maxAge"`
	RotateTime int    `json:"rotate_time" yaml:"rotateTime"`
}

type Config struct {
	FilePath         string `json:"file_path" yaml:"filePath"`
	FileName         string `json:"file_name" yaml:"fileName"`
	Level            string `json:"level" yaml:"level"`
	RotateType       string `json:"rotate_type" yaml:"rotateType"` // time size
	FormatType       string `json:"format_type" yaml:"formatType"` // json console
	Compress         bool   `json:"compress" yaml:"compress"`
	EnableStackTrace bool   // 是否开启调用栈全量信息

	SizeRotate *SizeRotate `json:"size_rotate" yaml:"sizeRotate"`
	TimeRotate *TimeRotate `json:"time_rotate" yaml:"timeRotate"`
}

func strToZapLevel(lv string) zapcore.Level {
	switch strings.ToLower(lv) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelDebug
	}
}
