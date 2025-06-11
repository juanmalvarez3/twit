package service

import (
	"context"
	"fmt"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"go.uber.org/zap"
)

func (s Service) Update(ctx context.Context, entry dmntimeline.TimelineEntry, userID string) error {
	s.logger.Debug("Obteniendo timeline",
		zap.String("action", actionUpdate),
		zap.String("user_id", userID))

	fmt.Println("Obteniendo timeline del repositorio con userID:", userID)
	err := s.timelineRepo.Update(ctx, entry, userID)
	if err != nil {
		s.logger.Error("Error al actualizar timeline",
			zap.String("action", actionUpdate),
			zap.String("user_id", userID),
			zap.Error(err))
		return err
	}

	s.logger.Debug("Timeline actualizado exitosamente",
		zap.String("action", actionUpdate),
		zap.String("user_id", userID))

	return nil
}
