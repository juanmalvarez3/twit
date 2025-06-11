package createfollow

import (
	"context"
	"fmt"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"go.uber.org/zap"
	"time"
)

func (u UseCase) CreateFollow(ctx context.Context, follow dmnfollow.Follow) error {
	if follow.FollowerID == follow.FollowedID {
		return fmt.Errorf("un usuario no puede seguirse a sí mismo")
	}

	if follow.CreatedAt == "" {
		follow.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	u.logger.Info("Creando nuevo follow",
		zap.String("follower_id", follow.FollowerID),
		zap.String("followed_id", follow.FollowedID))

	if err := u.service.Create(ctx, follow); err != nil {
		u.logger.Error("Error al persistir follow",
			zap.Error(err))
		return err
	}

	u.logger.Info("Follow creado exitosamente",
		zap.String("follower_id", follow.FollowerID),
		zap.String("followed_id", follow.FollowedID))

	return nil
}
