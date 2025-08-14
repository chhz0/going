package zlog

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogger_Info(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		fields   []Field
		wantLogs int
	}{
		{
			name:     "simple message",
			msg:      "test message",
			fields:   nil,
			wantLogs: 1,
		},
		{
			name:     "with fields",
			msg:      "test with fields",
			fields:   []Field{zap.String("key", "value")},
			wantLogs: 1,
		},
		{
			name:     "empty message",
			msg:      "",
			fields:   nil,
			wantLogs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup observer
			core, recorded := observer.New(zapcore.InfoLevel)
			level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
			logger := &zapLogger{
				l:   zap.New(core),
				lvl: &level,
			}

			// Call method
			logger.Info(tt.msg, tt.fields...)

			// Verify
			logs := recorded.All()
			if len(logs) != tt.wantLogs {
				t.Errorf("Expected %d logs, got %d", tt.wantLogs, len(logs))
			}

			if len(logs) > 0 {
				if logs[0].Message != tt.msg {
					t.Errorf("Expected message %q, got %q", tt.msg, logs[0].Message)
				}
			}
		})
	}
}

func TestZapLogger_Debug(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	level := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger := &zapLogger{
		l:   zap.New(core),
		lvl: &level,
	}

	logger.Debug("debug message", zap.String("key", "value"))

	logs := recorded.All()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	if logs[0].Message != "debug message" {
		t.Errorf("Expected message 'debug message', got %q", logs[0].Message)
	}

	if logs[0].Level != zapcore.DebugLevel {
		t.Errorf("Expected level DebugLevel, got %v", logs[0].Level)
	}
}

func TestZapLogger_Warn(t *testing.T) {
	core, recorded := observer.New(zapcore.WarnLevel)
	level := zap.NewAtomicLevelAt(zapcore.WarnLevel)
	logger := &zapLogger{
		l:   zap.New(core),
		lvl: &level,
	}

	logger.Warn("warning message", zap.String("key", "value"))

	logs := recorded.All()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	if logs[0].Message != "warning message" {
		t.Errorf("Expected message 'warning message', got %q", logs[0].Message)
	}

	if logs[0].Level != zapcore.WarnLevel {
		t.Errorf("Expected level WarnLevel, got %v", logs[0].Level)
	}
}

func TestZapLogger_Error(t *testing.T) {
	core, recorded := observer.New(zapcore.ErrorLevel)
	level := zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	logger := &zapLogger{
		l:   zap.New(core),
		lvl: &level,
	}

	logger.Error("error message", zap.String("key", "value"))

	logs := recorded.All()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	if logs[0].Message != "error message" {
		t.Errorf("Expected message 'error message', got %q", logs[0].Message)
	}

	if logs[0].Level != zapcore.ErrorLevel {
		t.Errorf("Expected level ErrorLevel, got %v", logs[0].Level)
	}
}

func TestZapLogger_WithName(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger := &zapLogger{
		l:   zap.New(core),
		lvl: &level,
	}

	namedLogger := logger.WithName("test-component")
	namedLogger.Info("test message")

	logs := recorded.All()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	entry := logs[0]
	if entry.LoggerName != "test-component" {
		t.Errorf("Expected logger name 'test-component', got %q", entry.LoggerName)
	}
}

func TestZapLogger_WithFields(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger := &zapLogger{
		l:   zap.New(core),
		lvl: &level,
	}

	fieldLogger := logger.WithFields(zap.String("version", "1.0"), zap.Int("count", 42))
	fieldLogger.Info("test message")

	logs := recorded.All()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	entry := logs[0]
	foundVersion := false
	foundCount := false
	for _, field := range entry.Context {
		if field.Key == "version" && field.String == "1.0" {
			foundVersion = true
		}
		if field.Key == "count" && field.Integer == 42 {
			foundCount = true
		}
	}

	if !foundVersion {
		t.Error("Expected field 'version' with value '1.0'")
	}

	if !foundCount {
		t.Error("Expected field 'count' with value 42")
	}
}

func TestZapLogger_WithContext(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger := &zapLogger{
		l:   zap.New(core),
		lvl: &level,
	}

	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	ctxLogger := logger.WithContext(ctx, "request_id")
	ctxLogger.Info("test message")

	logs := recorded.All()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	entry := logs[0]
	foundRequestID := false
	for _, field := range entry.Context {
		if field.Key == "request_id" && field.String == "req-123" {
			foundRequestID = true
		}
	}

	if !foundRequestID {
		t.Error("Expected field 'request_id' with value 'req-123'")
	}
}

func TestZapLogger_LevelControl(t *testing.T) {
	core, recorded := observer.New(zapcore.WarnLevel)
	level := zap.NewAtomicLevelAt(zapcore.WarnLevel)
	logger := &zapLogger{
		l:   zap.New(core),
		lvl: &level,
	}

	// Debug should not be logged
	logger.Debug("debug message")
	if logs := recorded.All(); len(logs) != 0 {
		t.Errorf("Expected 0 logs for debug message, got %d", len(logs))
	}

	// Warn should be logged
	logger.Warn("warn message")
	if logs := recorded.All(); len(logs) != 1 {
		t.Errorf("Expected 1 log for warn message, got %d", len(logs))
	}
}

func TestGlobalLoggerSingleton(t *testing.T) {
	// Store original logger
	originalLogger := std.Load()
	defer std.Store(originalLogger) // Restore original logger

	// Create test logger
	core, recorded := observer.New(zapcore.InfoLevel)
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	testLogger := &zapLogger{
		l:   zap.New(core),
		lvl: &level,
	}

	// Replace global logger
	std.Store(testLogger)

	// Test global functions
	Info("global info message")
	Debug("global debug message")
	Warn("global warn message")

	logs := recorded.All()
	if len(logs) != 2 { // Debug won't be recorded due to default level
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}

	// Check info message
	infoLogs := recorded.FilterLevelExact(zapcore.InfoLevel).TakeAll()
	if len(infoLogs) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(infoLogs))
	} else if infoLogs[0].Message != "global info message" {
		t.Errorf("Expected message 'global info message', got %q", infoLogs[0].Message)
	}

	// Check warn message
	warnLogs := recorded.FilterLevelExact(zapcore.WarnLevel).TakeAll()
	if len(warnLogs) != 1 {
		t.Errorf("Expected 1 warn log, got %d", len(warnLogs))
	} else if warnLogs[0].Message != "global warn message" {
		t.Errorf("Expected message 'global warn message', got %q", warnLogs[0].Message)
	}
}
