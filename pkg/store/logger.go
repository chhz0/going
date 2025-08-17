package store

import (
	"context"
)

// Logger 日志接口
type Logger interface {
	Error(ctx context.Context, msg string, fields ...interface{})
}

type emptyLogger struct{}

func (l emptyLogger) Error(_ context.Context, _ string, _ ...interface{}) {}

// New 创建空日志记录器
func New() Logger {
	return emptyLogger{}
}
