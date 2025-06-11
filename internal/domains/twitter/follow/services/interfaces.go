package services

import (
	"context"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain/events"
)

type Repository interface {
	Create(ctx context.Context, follow dmnfollow.Follow) error
	Get(ctx context.Context, followID string) (dmnfollow.Follow, error)
	GetFollowers(ctx context.Context, followedID string) ([]string, error)
	GetFollowing(ctx context.Context, followerID string) ([]string, error)
}

type Publisher interface {
	Publish(ctx context.Context, event events.Event) error
}
