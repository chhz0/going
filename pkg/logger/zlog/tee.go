package zlog

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LevelEnableFunc func(level Level) bool

type TeeOption struct {
	Output io.Writer
	LevelEnableFunc
}

// newTee creates a new tee logger.
// 默认使用的 日志级别为 InfoLevel，
func newTee(tees []TeeOption, encoder EncoderType, opts ...ZapLoggerOption) *zapLogger {
	cores := make([]zapcore.Core, 0, len(tees))

	for _, tee := range tees {
		if tee.Output == nil {
			tee.Output = os.Stdout
		}

		encoderConfig := &zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		core := zapcore.NewCore(
			switchEncoder(encoder, encoderConfig),
			zapcore.AddSync(tee.Output),
			zap.LevelEnablerFunc(tee.LevelEnableFunc),
		)
		cores = append(cores, core)
	}

	atomicl := zap.NewAtomicLevelAt(InfoLevel)
	return &zapLogger{
		l:   zap.New(zapcore.NewTee(cores...), opts...),
		lvl: &atomicl,
	}
}

func OpenLogFile(file string) io.Writer {
	logf, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("open log file failed")
	}
	return logf
}

func NewTeeLogger(tees []TeeOption, encoder EncoderType, opts ...ZapLoggerOption) Logger {
	return newTee(tees, encoder, opts...)
}
