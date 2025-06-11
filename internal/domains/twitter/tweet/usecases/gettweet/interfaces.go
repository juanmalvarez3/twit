package gettweet

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"go.uber.org/zap"
)

type TweetsService interface {
	Get(ctx context.Context, id string) (dmntweet.Tweet, error)
	Search(ctx context.Context, opts options.SearchOptions) ([]dmntweet.Tweet, string, error)
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}
