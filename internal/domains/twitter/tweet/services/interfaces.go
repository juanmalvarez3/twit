package services

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/events"
)

type Repository interface {
	Create(ctx context.Context, tweet dmntweet.Tweet) error
	Get(ctx context.Context, tweetID string) (dmntweet.Tweet, error)
	Search(ctx context.Context, userID string, limit int, lastEvaluatedKey string) ([]dmntweet.Tweet, string, error)
}

type Publisher interface {
	Publish(ctx context.Context, event events.Event) error
}
