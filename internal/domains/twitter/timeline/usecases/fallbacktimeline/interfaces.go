package fallbacktimeline

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	dmnoptions "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
)

type FollowerService interface {
	GetAllFollowing(ctx context.Context, followerID string) ([]string, error)
}

type TweetService interface {
	Search(ctx context.Context, opts dmnoptions.SearchOptions) ([]dmntweet.Tweet, string, error)
}

type Publisher interface {
	Publish(ctx context.Context, userID string, entries []dmntimeline.TimelineEntry) error
}

type TimelineService interface {
	BulkInsert(ctx context.Context, userID string, entries []dmntimeline.TimelineEntry) error
}

type CacheService interface {
	SetTimeline(ctx context.Context, timeline dmntimeline.Timeline) error
}
