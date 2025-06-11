package mocks

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type FollowerService struct {
	mock.Mock
}

func (m *FollowerService) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

type Publisher struct {
	mock.Mock
}

func (m *Publisher) Publish(ctx context.Context, tweet dmntweet.Tweet, timelineID string) error {
	args := m.Called(ctx, tweet, timelineID)
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
