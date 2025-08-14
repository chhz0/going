package zlog

import (
	"context"
	"io"
	"os"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

// ILogger 基础日志接口
type ILogger interface {
	Print(v ...any)
	Printf(format string, args ...any)

	Enable() bool
}

type zapILogger struct {
	sl *zap.SugaredLogger
}

func (l *zapILogger) Print(v ...any) {
	l.sl.Info(v...)
}

func (l *zapILogger) Printf(format string, args ...any) {
	l.sl.Infof(format, args...)
}

func (l *zapILogger) Enable() bool {
	return l.sl.Desugar().Core().Enabled(zap.InfoLevel)
}

// Logger 全日志业务接口
type Logger interface {
	Debug(msg string, field ...Field)
	Debugf(format string, args ...any)
	Debugw(msg string, keyAndValues ...any)
	Info(msg string, field ...Field)
	Infof(format string, args ...any)
	Infow(msg string, keyAndValues ...any)
	Warn(msg string, field ...Field)
	Warnf(format string, args ...any)
	Warnw(msg string, keyAndValues ...any)
	Error(msg string, field ...Field)
	Errorf(format string, args ...any)
	Errorw(msg string, keyAndValues ...any)
	Panic(msg string, field ...Field)
	Panicf(format string, args ...any)
	Panicw(msg string, keyAndValues ...any)
	Fatal(msg string, field ...Field)
	Fatalf(format string, args ...any)
	Fatalw(msg string, keyAndValues ...any)

	// WithName 为 Logger 附带名称
	WithName(name string) Logger
	// WithFields 为 Logger 附带结构化字段
	WithFields(fields ...Field) Logger
	// WithContext 为 Logger 附带 context 中的 key/values
	WithContext(ctx context.Context, keys ...string) Logger

	// V 设置动态日志级别, 但仅允许调用基础日志接口
	V(lvl Level) ILogger

	Sync()
}

type LoggerExtend interface {
	logger.Interface // gorm.logger
}

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

var std atomic.Pointer[zapLogger]

func init() {
	stdlog := New(
		os.Stdout,
		InfoLevel,
		JSONEncoder,
		WithCaller(true),
		Development(),
		ErrorOutput(zapcore.AddSync(os.Stderr)),
		AddStacktrace(zapcore.PanicLevel),
		AddCallerSkip(1),
		WithFatalHook(zapcore.WriteThenNoop),
	)
	std.Store(stdlog)
}

func ReplaceDefault(l Logger) {
	std.Store(l.(*zapLogger))
}

type zapLogger struct {
	l   *zap.Logger
	lvl *zap.AtomicLevel
}

func (l *zapLogger) Info(msg string, field ...Field) {
	l.l.Info(msg, field...)
}

func Info(msg string, field ...Field) {
	std.Load().Info(msg, field...)
}

func (l *zapLogger) Infof(format string, args ...any) {
	l.l.Sugar().Infof(format, args...)
}

func Infof(format string, args ...any) {
	std.Load().Infof(format, args...)
}

func (l *zapLogger) Infow(msg string, keyAndValues ...any) {
	l.l.Sugar().Infow(msg, keyAndValues...)
}

func Infow(msg string, keyAndValues ...any) {
	std.Load().Infow(msg, keyAndValues...)
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func Debug(msg string, fields ...Field) {
	std.Load().Debug(msg, fields...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.l.Sugar().Debugf(format, args...)
}

func Debugf(format string, args ...any) {
	std.Load().Debugf(format, args...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...any) {
	l.l.Sugar().Debugw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...any) {
	std.Load().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	std.Load().Warn(msg, fields...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.l.Sugar().Warnf(format, args...)
}

func Warnf(format string, args ...any) {
	std.Load().Warnf(format, args...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...any) {
	l.l.Sugar().Warnw(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...any) {
	std.Load().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func Error(msg string, fields ...Field) {
	std.Load().Error(msg, fields...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.l.Sugar().Errorf(format, args...)
}

func Errorf(format string, args ...any) {
	std.Load().Errorf(format, args...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...any) {
	l.l.Sugar().Errorw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...any) {
	std.Load().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	std.Load().Panic(msg, fields...)
}

func (l *zapLogger) Panicf(format string, args ...any) {
	l.l.Sugar().Panicf(format, args...)
}

func Panicf(format string, args ...any) {
	std.Load().Panicf(format, args...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...any) {
	l.l.Sugar().Panicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...any) {
	std.Load().Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	std.Load().Fatal(msg, fields...)
}

func (l *zapLogger) Fatalf(format string, args ...any) {
	l.l.Sugar().Fatalf(format, args...)
}

func Fatalf(format string, args ...any) {
	std.Load().Fatalf(format, args...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...any) {
	l.l.Sugar().Fatalw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...any) {
	std.Load().Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) WithName(name string) Logger {
	namedLogger := l.l.Named(name)
	return &zapLogger{
		l:   namedLogger,
		lvl: l.lvl,
	}
}

func WithName(name string) Logger {
	return std.Load().WithName(name)
}

func (l *zapLogger) WithFields(fields ...Field) Logger {
	newLogger := l.l
	if len(fields) > 0 {
		newLogger = l.l.With(fields...)
	}
	return &zapLogger{
		l:   newLogger,
		lvl: l.lvl,
	}
}

func WithFields(fields ...Field) Logger {
	return std.Load().WithFields(fields...)
}

func (l *zapLogger) WithContext(ctx context.Context, keys ...string) Logger {
	if len(keys) == 0 || ctx == nil {
		return l
	}

	kvs := make(map[string]any, len(keys))
	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			kvs[key] = val
		}
	}

	newLogger := l.clone()
	fields := make([]Field, 0, len(keys))
	for k, v := range kvs {
		fields = append(fields, AnyField(k, v))
	}

	if len(fields) > 0 {
		newLogger.l = newLogger.l.With(fields...)
	}
	return newLogger
}

func WithContext(ctx context.Context, keys ...string) Logger {
	return std.Load().WithContext(ctx, keys...)
}

func (l *zapLogger) V(lvl Level) ILogger {
	newLogger := l.clone()
	if newLogger.lvl != nil {
		newLogger.lvl.SetLevel(lvl)
	}
	return &zapILogger{
		sl: newLogger.l.Sugar(),
	}
}

func V(lvl Level) ILogger {
	return std.Load().V(lvl)
}

func (l *zapLogger) Sync() {
	if err := l.l.Sync(); err != nil {
		l.l.Error("failed to sync logger", zap.Error(err))
	}
}

func (l *zapLogger) SyncWithHandler(handler func(error)) {
	if err := l.l.Sync(); err != nil && handler != nil {
		handler(err)
	}
}

func Sync() {
	std.Load().Sync()
}

func SyncWithHandler(handler func(error)) {
	std.Load().SyncWithHandler(handler)
}

func (l *zapLogger) SetLevel(lvl Level) {
	if l.lvl != nil {
		l.lvl.SetLevel(lvl)
	}
}

func SetLevel(lvl Level) {
	std.Load().SetLevel(lvl)
}

func (l *zapLogger) Zap() *zap.Logger {
	return l.l
}

func Zap() *zap.Logger {
	return std.Load().Zap()
}

func (l *zapLogger) clone() *zapLogger {
	clone := *l
	return &clone
}

type Output func() io.Writer

type EncoderType string

const (
	ConsoleEncoder EncoderType = "console"
	JSONEncoder    EncoderType = "json"
)

func New(out io.Writer, lvl Level, encoder EncoderType, opts ...ZapLoggerOption) *zapLogger {
	core, al := newCore(out, lvl, encoder)
	return &zapLogger{
		l:   zap.New(core, opts...),
		lvl: al,
	}
}

func newCore(out io.Writer, lvl Level, encoder EncoderType) (zapcore.Core, *zap.AtomicLevel) {
	if out == nil {
		out = os.Stdout
	}

	atomicl := zap.NewAtomicLevelAt(lvl)
	enc := &zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewCore(
		switchEncoder(encoder, enc),
		zapcore.AddSync(out),
		atomicl,
	), &atomicl
}

func switchEncoder(encoder EncoderType, encoderConfig *zapcore.EncoderConfig) zapcore.Encoder {
	switch encoder {
	case ConsoleEncoder:
		return zapcore.NewConsoleEncoder(*encoderConfig)
	case JSONEncoder:
		return zapcore.NewJSONEncoder(*encoderConfig)
	default:
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
}
