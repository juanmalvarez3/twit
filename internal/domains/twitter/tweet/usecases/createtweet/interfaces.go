package createtweet

import (
	"context"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"go.uber.org/zap"
)

type TweetsService interface {
	Get(ctx context.Context, id string) (domain.Tweet, error)
	Create(ctx context.Context, twt domain.Tweet) (domain.Tweet, error)
	Search(ctx context.Context, opts options.SearchOptions) ([]domain.Tweet, string, error)
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}
