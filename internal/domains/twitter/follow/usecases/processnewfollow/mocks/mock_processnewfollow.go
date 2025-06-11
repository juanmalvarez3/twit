package mocks

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
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

type TweetsService struct {
	mock.Mock
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

type UpdateTimelineSQSPublisher struct {
	mock.Mock
}

func (m *UpdateTimelineSQSPublisher) Publish(ctx context.Context, tweet dmntweet.Tweet, userID string) error {
	args := m.Called(ctx, tweet, userID)
	return args.Error(0)
}
