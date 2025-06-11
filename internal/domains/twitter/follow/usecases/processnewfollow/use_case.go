package processnewfollow

import (
	"context"
	"go.uber.org/zap"
)

type UseCase struct {
	service                 Service
	logger                  Logger
	tweetService            TweetsService
	updateTimelinePublisher UpdateTimelineSQSPublisher
}

func NewUseCase(
	service Service,
	logger Logger,
	tweetsService TweetsService,
	updateTimelinePublisher UpdateTimelineSQSPublisher,
) UseCase {
	return UseCase{
		service:                 service,
		logger:                  logger,
		tweetService:            tweetsService,
		updateTimelinePublisher: updateTimelinePublisher,
	}
}

func (uc *UseCase) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	uc.logger.Info("Obteniendo seguidores", zap.String("user_id", userID))

	followers, err := uc.service.GetFollowers(ctx, userID)
	if err != nil {
		uc.logger.Error("Error al obtener seguidores",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	uc.logger.Debug("Seguidores obtenidos exitosamente",
		zap.String("user_id", userID),
		zap.Int("count", len(followers)))

	return followers, nil
}

func (uc *UseCase) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	uc.logger.Info("Obteniendo usuarios seguidos", zap.String("user_id", userID))

	following, err := uc.service.GetFollowing(ctx, userID)
	if err != nil {
		uc.logger.Error("Error al obtener usuarios seguidos",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	uc.logger.Debug("Usuarios seguidos obtenidos exitosamente",
		zap.String("user_id", userID),
		zap.Int("count", len(following)))

	return following, nil
}
