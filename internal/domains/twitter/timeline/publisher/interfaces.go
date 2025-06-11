package publisher

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
)

type TimelinePublisher interface {
	Publish(ctx context.Context, timeline dmntimeline.Timeline) error
}

type RebuildPublisher interface {
	Publish(ctx context.Context, userID string) error
}
