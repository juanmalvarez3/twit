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

func (m *TimelineService) Get(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, bool, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).(dmntimeline.Timeline), args.Bool(1), args.Error(2)
}

type Publisher struct {
	mock.Mock
}

func (m *Publisher) Publish(ctx context.Context, timeline dmntimeline.Timeline) error {
	args := m.Called(ctx, timeline)
	return args.Error(0)
}

type FallbackRebuildTimelinePublisherService struct {
	mock.Mock
}

func (m *FallbackRebuildTimelinePublisherService) Publish(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
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
