package orchestratefanout

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"go.uber.org/zap"
)

func (u *UseCase) Exec(ctx context.Context, tweet dmntweet.Tweet) error {
	u.logger.Info("Iniciando distribución de tweet a timelines de seguidores",
		zap.String("tweet_id", tweet.ID),
		zap.String("user_id", tweet.UserID))

	followers, err := u.followerService.GetFollowers(ctx, tweet.UserID)
	if err != nil {
		u.logger.Error("Error obteniendo seguidores para distribución de tweet",
			zap.String("tweet_id", tweet.ID),
			zap.String("user_id", tweet.UserID),
			zap.Error(err))
		return err
	}

	u.logger.Debug("Seguidores obtenidos para distribución",
		zap.String("tweet_id", tweet.ID),
		zap.String("user_id", tweet.UserID),
		zap.Int("followers_count", len(followers)))

	if len(followers) == 0 {
		u.logger.Debug("No hay seguidores para distribuir el tweet",
			zap.String("tweet_id", tweet.ID),
			zap.String("user_id", tweet.UserID))
		return nil
	}

	for _, followerID := range followers {
		tweetCopy := dmntweet.Tweet{
			ID:        tweet.ID,
			UserID:    tweet.UserID,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt,
		}

		err := u.publisher.Publish(ctx, tweetCopy, followerID)
		if err != nil {
			u.logger.Error("Error publicando evento de actualización de timeline",
				zap.String("tweet_id", tweet.ID),
				zap.String("user_id", tweet.UserID),
				zap.String("follower_id", followerID),
				zap.Error(err))
			continue
		}

		u.logger.Debug("Evento de actualización de timeline publicado",
			zap.String("tweet_id", tweet.ID),
			zap.String("user_id", tweet.UserID),
			zap.String("follower_id", followerID))
	}

	u.logger.Info("Distribución de tweet completada",
		zap.String("tweet_id", tweet.ID),
		zap.String("user_id", tweet.UserID),
		zap.Int("followers_processed", len(followers)))

	return nil
}
