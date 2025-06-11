package processnewfollow

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"go.uber.org/zap"
)

type Service interface {
	GetFollowers(ctx context.Context, followedID string) ([]string, error)
	GetFollowing(ctx context.Context, followerID string) ([]string, error)
}

type TweetsService interface {
	Search(ctx context.Context, opts options.SearchOptions) ([]dmntweet.Tweet, string, error)
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type UpdateTimelineSQSPublisher interface {
	Publish(ctx context.Context, tweet dmntweet.Tweet, userID string) error
}
