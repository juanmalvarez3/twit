package createfollow

import (
	"context"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, follow dmnfollow.Follow) error
}

type Publisher interface {
	PublishFollowCreated(ctx context.Context, follow dmnfollow.Follow) error
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}
