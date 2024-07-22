package ctxlog

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	trackLogger *CtxLogger
)

// track log 扩展信息,用于非必填项 见日志规范 track日志
type ExtTrackInfo struct {
	ErrMsg   string `json:"errmsg"`   // 本次请求服务对外提供的返回短语
	Req      string `json:"req"`      // 本次请求的原始入参
	Resp     string `json:"resp"`     // 本次请求的返回结果
	DeviceID string `json:"deviceid"` // 从客户端带过来的设备ID
}

// NewTrackLogger 生成track日志logger
// @serviceName 服务/模块名称
// @filePath 日志所在目录
// @maxAge 日志备份最长时间 单位小时
// @rotateTime 多久切分一次 单位小时
// @ctxFunc 见bfo_ctx_logger type CtxFunc func(ctx context.Context) map[string]any, 一个解析ctx中数据的函数,用于日志记录
func NewTrackLogger(serviceName string, filePath string, maxAge int, rotateTime int, exposedKey []string) *CtxLogger {
	hook := timeRotateHook(filePath+"/track.log", &TimeRotateConfig{
		FileNameFormat: "%Y-%m-%d",
		MaxAge:         maxAge,
		RotateTime:     rotateTime,
	})
	encoderConfig := zapcore.EncoderConfig{
		TimeKey: "time",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
		},
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.InfoLevel)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(hook),
		atomicLevel,
	)
	sugar := zap.New(core, zap.Fields(zap.String("svr", serviceName))).Sugar()
	// 单独new track logger, default和track默认用同一个logger
	return NewCtxLogger(sugar, sugar, exposedKey)

}

// InitTrackLogger 默认初始化track logger
// @serviceName 服务/模块名称
// @ctxFunc 见bfo_ctx_logger type CtxFunc func(ctx context.Context) map[string]any, 一个解析ctx中数据的函数,用于日志记录
func InitTrackLogger(serviceName string, exposedKey []string) {
	logger := NewTrackLogger(serviceName, "./logs/", 24*7, 24, exposedKey)
	SetTrackLogger(logger)
}

// InitStdoutTrackLogger 以stdout作为输出的track logger，用于debug和docker部署
// @serviceName 服务/模块名称
// @ctxFunc 见bfo_ctx_logger type CtxFunc func(ctx context.Context) map[string]any, 一个解析ctx中数据的函数,用于日志记录
func InitStdoutTrackLogger(serviceName string, exposedKey []string) {
	logger := NewStdoutLoggerWithExposedKey(serviceName, "track", JsonFormat, "info", exposedKey)
	SetTrackLogger(logger)
}

func SetTrackLogger(lg *CtxLogger) {
	trackLogger = lg
	// 程序退出时，通过调用logger.Sync()输出buffer中的内容
	go func() {
		c := make(chan os.Signal, 5)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		for {
			sig := <-c
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				signal.Stop(c)
				trackLogger.Sync()
			default:
				continue
			}
		}
	}()
}

func GetTrackLogger() *CtxLogger {
	return trackLogger
}
