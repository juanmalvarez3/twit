package services

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"go.uber.org/zap"
)

func (s Service) Get(ctx context.Context, id string) (dmntweet.Tweet, error) {
	s.logger.Debug("Obteniendo tweet", 
		zap.String("tweet_id", id),
		zap.String("action", actionGet),
	)

	twt, err := s.repository.Get(ctx, id)
	if err != nil {
		s.logger.Error("Error al obtener tweet", 
			zap.String("tweet_id", id),
			zap.Error(err),
			zap.String("action", actionGet),
		)
		return dmntweet.Tweet{}, err
	}

	if twt.ID == "" {
		s.logger.Warn("Tweet no encontrado", 
			zap.String("tweet_id", id),
			zap.String("action", actionGet),
		)
		return dmntweet.Tweet{}, dmntweet.ErrTweetNotFound
	}

	s.logger.Debug("Tweet obtenido exitosamente", 
		zap.String("tweet_id", twt.ID),
		zap.String("user_id", twt.UserID),
		zap.String("action", actionGet),
	)

	return twt, nil
}
