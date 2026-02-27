package xcontext

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	startTimeKey = "_xcontext_start_time_"
	traceKey     = "_xcontext_trace_"
)

func WithTrace(ctx context.Context, trace string) context.Context {
	return context.WithValue(ctx, traceKey, trace)
}

func WithStartTime(ctx context.Context, start time.Time) context.Context {
	return context.WithValue(ctx, startTimeKey, start)
}

func Background() context.Context {
	ctx := context.Background()
	ctx = WithStartTime(ctx, time.Now())
	ctx = WithTrace(ctx, uuid.New().String())
	return ctx
}

func GetTrace(ctx context.Context) string {
	temp := ctx.Value(traceKey)
	trace, ok := temp.(string)
	if !ok {
		return "[unknown trace]"
	}
	return trace
}

func Since(ctx context.Context) time.Duration {
	temp := ctx.Value(startTimeKey)
	start, ok := temp.(time.Time)
	if !ok {
		return 0
	}
	return time.Since(start)
}
