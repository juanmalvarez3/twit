package mocks

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	dmnoptions "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"github.com/stretchr/testify/mock"
)

type FollowerService struct {
	mock.Mock
}

func (m *FollowerService) GetAllFollowing(ctx context.Context, followerID string) ([]string, error) {
	args := m.Called(ctx, followerID)
	return args.Get(0).([]string), args.Error(1)
}

type TweetService struct {
	mock.Mock
}

func (m *TweetService) Search(ctx context.Context, opts dmnoptions.SearchOptions) ([]dmntweet.Tweet, string, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]dmntweet.Tweet), args.String(1), args.Error(2)
}

type Publisher struct {
	mock.Mock
}

func (m *Publisher) Publish(ctx context.Context, userID string, entries []dmntimeline.TimelineEntry) error {
	args := m.Called(ctx, userID, entries)
	return args.Error(0)
}

type TimelineService struct {
	mock.Mock
}

func (m *TimelineService) BulkInsert(ctx context.Context, userID string, entries []dmntimeline.TimelineEntry) error {
	args := m.Called(ctx, userID, entries)
	return args.Error(0)
}

type CacheService struct {
	mock.Mock
}

func (m *CacheService) SetTimeline(ctx context.Context, timeline dmntimeline.Timeline) error {
	args := m.Called(ctx, timeline)
	return args.Error(0)
}
