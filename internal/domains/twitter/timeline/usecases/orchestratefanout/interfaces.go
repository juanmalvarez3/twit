package orchestratefanout

import (
	"context"

	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"go.uber.org/zap"
)

type FollowerService interface {
	GetFollowers(ctx context.Context, userID string) ([]string, error)
}

type Publisher interface {
	Publish(ctx context.Context, tweet dmntweet.Tweet, timelineID string) error
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}
