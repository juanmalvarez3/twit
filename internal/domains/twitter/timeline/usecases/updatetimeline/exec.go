package updatetimeline

import (
	"context"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"go.uber.org/zap"
)

func (u *UseCase) Exec(ctx context.Context, tweet dmntweet.Tweet, userID string) error {
	u.logger.Debug("Tweet recibido para actualizar timeline",
		zap.String("tweet_id", tweet.ID),
		zap.String("user_id", tweet.UserID),
		zap.String("content", tweet.Content),
		zap.String("created_at", tweet.CreatedAt))

	entry := dmntimeline.NewTimelineEntryFromTweet(tweet)

	u.logger.Debug("Entrada de timeline creada",
		zap.String("tweet_id", entry.TweetID),
		zap.String("author_id", entry.AuthorID),
		zap.String("content", entry.Content),
		zap.Time("created_at", entry.CreatedAt))

	return u.timelineService.Update(ctx, entry, userID)
}
