package ctxlog

import (
	"context"
	"time"
)

const CTX_KEY = "CONTEXT_DATA_KEY"

type Ctx struct {
	UUID      string         // 本次请求的唯一id
	Infc      string         // 接口名称
	Ext       map[string]any // 用于自定义set一些值
	StartTime time.Time      // 开始处理请求的 starttime
}

// GenDefaultCtxData 生成默认的atlantis ctx data 会携带一些默认的环境信息
func GenDefaultCtxData(ctx context.Context) (data Ctx) {
	return Ctx{}
}

func DefaultCtx(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	oldCtx, ok := GetCtxData(ctx)
	newCtx := GenDefaultCtxData(ctx)
	if ok {
		return WithCtxData(ctx, oldCtx)
	}
	return WithCtxData(ctx, newCtx)
}

func GetCtxData(ctx context.Context) (Ctx, bool) {
	if ctx == nil {
		return Ctx{}, false
	}
	data, exist := ctx.Value(CTX_KEY).(Ctx)
	if !exist {
		data = Ctx{}
	}
	return data, exist
}

func WithCtxData(ctx context.Context, data Ctx) context.Context {
	return context.WithValue(ctx, CTX_KEY, data)
}
