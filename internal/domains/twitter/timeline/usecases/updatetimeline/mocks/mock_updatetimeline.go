package mocks

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type TimelineService struct {
	mock.Mock
}

func (m *TimelineService) Update(ctx context.Context, entry dmntimeline.TimelineEntry, userID string) error {
	args := m.Called(ctx, entry, userID)
	return args.Error(0)
}

type Logger struct {
	mock.Mock
}

func (m *Logger) Debug(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *Logger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *Logger) Warn(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *Logger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *Logger) Fatal(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *Logger) With(fields ...zap.Field) *Logger {
	args := m.Called(fields)
	return args.Get(0).(*Logger)
}

func (m *Logger) WithValues(keysAndValues ...interface{}) *Logger {
	args := m.Called(keysAndValues)
	return args.Get(0).(*Logger)
}

func (m *Logger) Named(name string) *Logger {
	args := m.Called(name)
	return args.Get(0).(*Logger)
}

func (m *Logger) Sync() error {
	args := m.Called()
	return args.Error(0)
}
