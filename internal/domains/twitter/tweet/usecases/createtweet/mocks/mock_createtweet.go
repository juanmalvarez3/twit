package mocks

import (
	"context"
	domain "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type TweetsService struct {
	mock.Mock
}

func (m *TweetsService) Get(ctx context.Context, id string) (domain.Tweet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Tweet), args.Error(1)
}

func (m *TweetsService) Create(ctx context.Context, twt domain.Tweet) (domain.Tweet, error) {
	args := m.Called(ctx, twt)
	return args.Get(0).(domain.Tweet), args.Error(1)
}

func (m *TweetsService) Search(ctx context.Context, opts options.SearchOptions) ([]domain.Tweet, string, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]domain.Tweet), args.String(1), args.Error(2)
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
