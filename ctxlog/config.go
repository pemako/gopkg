package ctxlog

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug = zapcore.DebugLevel
	LevelInfo  = zapcore.InfoLevel
	LevelWarn  = zapcore.WarnLevel
	LevelError = zapcore.ErrorLevel
	LevelTrack = zapcore.InfoLevel
)

const (
	JsonFormat    = "json"
	ConsoleFormat = "console"

	RotateTypeTime = "time"
	RotateTypeSize = "size"
)

type SizeRotateConfig struct {
	MaxSize    int  `json:"max_size" yaml:"maxSize"`       // 每个日志文件保存的最大尺寸，单位M
	MaxBackups int  `json:"max_backups" yaml:"maxBackups"` // 每个日志文件保存的最大尺寸，单位M
	MaxAge     int  `json:"max_age" yaml:"maxAge"`         // 日志文件最多保存多少天
	Compress   bool `json:"compress" yaml:"compress"`      // 是否压缩
}

type TimeRotateConfig struct {
	FileNameFormat string `json:"format" yaml:"format"`          // 日志文件名称格式 如:"%Y%m%d%H%M"
	MaxAge         int    `json:"max_age" yaml:"maxAge"`         // 日志文件名称格式 如:"%Y%m%d%H%M"
	RotateTime     int    `json:"rotate_time" yaml:"rotateTime"` // 多少小时切割一次日志
}

type Config struct {
	ServiceName      string `json:"service_name" yaml:"serviceName"`            // 服务名称 必填
	FilePath         string `json:"file_path" yaml:"filePath"`                  // FileNameFormat
	Level            string `json:"level" yaml:"level"`                         // 日志文件切割方式
	FileName         string `json:"file_name" yaml:"fileName"`                  // 文件名 如果为空则使用 serviceName
	RotateType       string `json:"rotate_type" yaml:"rotateType"`              // 日志文件切割方式 time size
	FormatType       string `json:"format_type" yaml:"formatType"`              // 日志格式 json console
	EnableStackTrace bool   `json:"enable_stack_trace" yaml:"EnableStackTrace"` // 是否开启调用栈全量信息
	EnableDev        bool   // 是否开启dev模式(for dpanic)

	SizeRotateConfig *SizeRotateConfig `json:"size_rotate" yaml:"sizeRotate"`
	TimeRotateConfig *TimeRotateConfig `json:"time_rotate" yaml:"timeRotate"`

	ExecDir string // 工程根目录，用于解析atlantis隐藏文件
}

func (c *Config) GetFileName() string {
	if c.FileName != "" {
		return c.FileName
	}
	return c.ServiceName
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
