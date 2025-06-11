package createtweet

import (
	"context"
	"github.com/google/uuid"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"go.uber.org/zap"
	"time"
)

func (u UseCase) CreateTweet(ctx context.Context, tweet *dmntweet.Tweet) (*dmntweet.Tweet, error) {
	u.logger.Debug("Validando tweet",
		zap.String("user_id", tweet.UserID),
		zap.Int("content_length", len(tweet.Content)),
	)

	err := tweet.Validate()
	if err != nil {
		u.logger.Error("Error de validación del tweet",
			zap.String("user_id", tweet.UserID),
			zap.String("content", tweet.Content),
			zap.Error(err),
		)
		return nil, err
	}

	tweet.Content = tweet.NormalizeContent()
	u.logger.Debug("Contenido normalizado",
		zap.String("user_id", tweet.UserID),
		zap.String("content", tweet.Content),
	)
	tweet.ID = "twt-" + uuid.New().String()
	tweet.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	created, err := u.twtService.Create(ctx, *tweet)
	if err != nil {
		u.logger.Error("Error al crear tweet en el servicio",
			zap.String("user_id", tweet.UserID),
			zap.String("content", tweet.Content),
			zap.Error(err),
		)
		return nil, err
	}

	if created.ID != "" {
		u.logger.Info("Tweet creado exitosamente",
			zap.String("tweet_id", created.ID),
			zap.String("user_id", created.UserID),
			zap.String("content", created.Content),
			zap.String("created_at", created.CreatedAt),
		)
	} else {
		u.logger.Warn("Tweet creado sin ID",
			zap.String("user_id", created.UserID),
		)
	}
	return &created, nil
}
