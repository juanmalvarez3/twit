package gettimeline

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"go.uber.org/zap"
)

type TimelineService interface {
	Get(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, bool, error)
}

type Publisher interface {
	Publish(ctx context.Context, timeline dmntimeline.Timeline) error
}
type FallbackRebuildTimelinePublisherService interface {
	Publish(ctx context.Context, userID string) error
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}
