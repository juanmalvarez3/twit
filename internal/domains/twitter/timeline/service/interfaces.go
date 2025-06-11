package service

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	evttweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/events"
)

type Repository interface {
	Update(ctx context.Context, entry dmntimeline.TimelineEntry, userID string) error
	Get(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, bool, error)
	GetFromDB(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, error)
	SetCache(ctx context.Context, key string, value []byte) error
}

type Publisher interface {
	Publish(ctx context.Context, twt dmntweet.Tweet, eventType evttweet.EventType) error
}
