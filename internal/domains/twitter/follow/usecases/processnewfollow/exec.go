package processnewfollow

import (
	"context"
	"fmt"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain/events"
	dmnoptions "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"go.uber.org/zap"
)

func (uc *UseCase) ProcessNewFollow(ctx context.Context, followEvent events.FollowCreatedEvent) error {
	uc.logger.Info("Procesando nuevo follow",
		zap.String("follower_id", followEvent.Follow.FollowerID),
		zap.String("followed_id", followEvent.Follow.FollowedID))

	tweets, _, err := uc.tweetService.Search(ctx, dmnoptions.SearchOptions{
		Filters:    dmnoptions.SearchFilters{UserID: &followEvent.Follow.FollowedID},
		Pagination: dmnoptions.SearchPagination{Limit: 10, Offset: 0}})

	if err != nil {
		uc.logger.Error("Error al buscar tweets",
			zap.String("followed_id", followEvent.Follow.FollowedID),
			zap.Error(err))
		return fmt.Errorf("error al buscar tweets: %w", err)
	}

	if uc.updateTimelinePublisher != nil {
		for _, tweet := range tweets {
			if err := uc.updateTimelinePublisher.Publish(ctx, tweet, followEvent.Follow.FollowedID); err != nil {
				uc.logger.Error("Error al publicar mensaje de actualización de timeline",
					zap.String("follower_id", followEvent.Follow.FollowerID),
					zap.String("followed_id", followEvent.Follow.FollowedID),
					zap.Error(err))
				return fmt.Errorf("error al publicar mensaje: %w", err)
			}
			uc.logger.Debug("Mensaje enviado a cola de actualización de timeline",
				zap.String("follower_id", followEvent.Follow.FollowerID),
				zap.String("followed_id", followEvent.Follow.FollowedID))
		}
	} else {
		uc.logger.Warn("No se ha configurado publicador para actualización de timeline",
			zap.String("follower_id", followEvent.Follow.FollowerID),
			zap.String("followed_id", followEvent.Follow.FollowedID))
	}
	return nil
}
