package services

import (
	"context"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain/events"
	"go.uber.org/zap"
)

func (s Service) Create(ctx context.Context, follow dmnfollow.Follow) error {
	s.logger.Debug("Servicio: Creando follow",
		zap.String("follow_id", follow.ID),
		zap.String("follower_id", follow.FollowerID),
		zap.String("followed_id", follow.FollowedID))

	err := s.repository.Create(ctx, follow)
	if err != nil {
		s.logger.Error("Error al crear follow",
			zap.String("follow_id", follow.ID),
			zap.String("follower_id", follow.FollowerID),
			zap.String("followed_id", follow.FollowedID),
			zap.Error(err))
		return err
	}

	err = s.publisher.Publish(ctx, events.Event{
		Type:   events.FollowCreatedEventType,
		Follow: follow,
	})

	if err != nil {
		s.logger.Error("Error publicando evento de follow creado",
			zap.String("follow_id", follow.ID),
			zap.String("follower_id", follow.FollowerID),
			zap.String("followed_id", follow.FollowedID),
			zap.Error(err))
		return err
	}
	return nil
}
