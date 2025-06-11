package updatetimeline

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
)

type TimelineService interface {
	Update(ctx context.Context, entry dmntimeline.TimelineEntry, userID string) error
}
