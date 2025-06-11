package mocks

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type TweetsService struct {
	mock.Mock
}

func (m *TweetsService) Get(ctx context.Context, id string) (dmntweet.Tweet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dmntweet.Tweet), args.Error(1)
}

func (m *TweetsService) Search(ctx context.Context, opts options.SearchOptions) ([]dmntweet.Tweet, string, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]dmntweet.Tweet), args.String(1), args.Error(2)
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
