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

func (m *TimelineService) GetFromDB(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).(dmntimeline.Timeline), args.Error(1)
}

func (m *TimelineService) UpdateCache(ctx context.Context, timeline dmntimeline.Timeline) error {
	args := m.Called(ctx, timeline)
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
