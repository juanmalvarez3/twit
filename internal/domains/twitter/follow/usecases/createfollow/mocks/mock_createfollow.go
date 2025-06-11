package mocks

import (
	"context"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type Service struct {
	mock.Mock
}

func (m *Service) Create(ctx context.Context, follow dmnfollow.Follow) error {
	args := m.Called(ctx, follow)
	return args.Error(0)
}

type Publisher struct {
	mock.Mock
}

func (m *Publisher) PublishFollowCreated(ctx context.Context, follow dmnfollow.Follow) error {
	args := m.Called(ctx, follow)
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
