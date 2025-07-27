package test

import (
	"context"
	"go.uber.org/zap"
)

type MockLogger struct {
	InfoFunc  func(ctx context.Context, msg string, fields ...zap.Field)
	WarnFunc  func(ctx context.Context, msg string, fields ...zap.Field)
	ErrorFunc func(ctx context.Context, msg string, fields ...zap.Field)
}

func (m *MockLogger) Info(ctx context.Context, msg string, fields ...zap.Field)  { if m.InfoFunc != nil { m.InfoFunc(ctx, msg, fields...) } }
func (m *MockLogger) Warn(ctx context.Context, msg string, fields ...zap.Field)  { if m.WarnFunc != nil { m.WarnFunc(ctx, msg, fields...) } }
func (m *MockLogger) Error(ctx context.Context, msg string, fields ...zap.Field) { if m.ErrorFunc != nil { m.ErrorFunc(ctx, msg, fields...) } } 