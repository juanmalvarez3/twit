package createfollow

import (
	"context"
	"fmt"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (u UseCase) CreateFollow(ctx context.Context, follow dmnfollow.Follow) error {
	if follow.FollowerID == follow.FollowedID {
		return fmt.Errorf("un usuario no puede seguirse a sí mismo")
	}

	if follow.CreatedAt == "" {
		follow.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	follow.ID = "flw-" + follow.FollowerID + "-" + follow.FollowedID
	followFromDB, err := u.service.Get(ctx, follow.ID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		u.logger.Error(fmt.Sprintf("Error al obtener follow con ID %s. Error: %s", follow.ID, err.Error()))
		return err
	}

	if followFromDB.ID != "" {
		err := fmt.Errorf("el follow ya existe")
		u.logger.Error("Error al crear follow", zap.Error(err))
		return err
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
