package gettweet

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"go.uber.org/zap"
)

func (u UseCase) GetTweet(ctx context.Context, tweetID string) (dmntweet.Tweet, error) {
	u.logger.Debug("Obteniendo tweet",
		zap.String("tweet_id", tweetID),
	)

	twt, err := u.twtService.Get(ctx, tweetID)
	if err != nil {
		u.logger.Error("Error al obtener tweet",
			zap.String("tweet_id", tweetID),
			zap.Error(err),
		)
		return dmntweet.Tweet{}, err
	}

	u.logger.Debug("Tweet obtenido exitosamente",
		zap.String("tweet_id", twt.ID),
		zap.String("user_id", twt.UserID),
	)

	return twt, nil
}

func (u UseCase) validateTweet(tweet dmntweet.Tweet) error {
	u.logger.Debug("Validando tweet",
		zap.String("tweet_id", tweet.ID),
		zap.String("user_id", tweet.UserID),
	)

	err := tweet.Validate()
	if err != nil {
		u.logger.Warn("Error validando tweet",
			zap.String("tweet_id", tweet.ID),
			zap.String("user_id", tweet.UserID),
			zap.Error(err),
		)
		return err
	}

	u.logger.Debug("Tweet validado correctamente",
		zap.String("tweet_id", tweet.ID),
	)
	return nil
}
