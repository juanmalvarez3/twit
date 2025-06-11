package services

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/events"
	"go.uber.org/zap"
)

func (s Service) Create(ctx context.Context, twt dmntweet.Tweet) (dmntweet.Tweet, error) {
	err := s.repository.Create(ctx, twt)
	if err != nil {
		return dmntweet.Tweet{}, err
	}

	err = s.publisher.Publish(ctx, events.Event{
		Type:  events.TweetCreatedEventType,
		Tweet: twt,
	})
	if err != nil {
		s.logger.Error("Falló publicar evento SNS",
			zap.String("tweet_id", twt.ID),
			zap.Error(err),
			zap.String("action", actionCreate),
		)
	}

	return twt, nil
}
