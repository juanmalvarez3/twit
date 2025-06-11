package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type Service struct {
	mock.Mock
}

func (m *Service) GetFollowers(ctx context.Context, followedID string) ([]string, error) {
	args := m.Called(ctx, followedID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *Service) GetFollowing(ctx context.Context, followerID string) ([]string, error) {
	args := m.Called(ctx, followerID)
	return args.Get(0).([]string), args.Error(1)
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
