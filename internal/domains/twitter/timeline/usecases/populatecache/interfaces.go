package populatecache

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"go.uber.org/zap"
)

type TimelineService interface {
	GetFromDB(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, error)
	UpdateCache(ctx context.Context, timeline dmntimeline.Timeline) error
}

type CacheService interface {
	SetTimeline(ctx context.Context, timeline dmntimeline.Timeline) error
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}
