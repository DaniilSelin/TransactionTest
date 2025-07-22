package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerKey = "logger"
	RequestID = "RequestID"
)

type Logger struct {
	l *zap.Logger
}

func New(cfgLog *zap.Config) (Logger, error) {
	// Добавляем энкодер времени и уровня вручную, т.к это функции
	cfgLog.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfgLog.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder 

	logger, err := cfgLog.Build()
	if err != nil {
		return Logger{}, fmt.Errorf("FATAL: failed to create logger: %w", err)
	}

	return Logger{l: logger}, nil
}

// addRequestID достает из контекста RequestID и добавляет его в поля
func (l *Logger) addRequestID(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	return fields
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Info(msg, l.addRequestID(ctx, fields)...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Debug(msg, l.addRequestID(ctx, fields)...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Warn(msg, l.addRequestID(ctx, fields)...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Error(msg, l.addRequestID(ctx, fields)...)
}

func (l *Logger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.DPanic(msg, l.addRequestID(ctx, fields)...)
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Panic(msg, l.addRequestID(ctx, fields)...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Fatal(msg, l.addRequestID(ctx, fields)...)
}