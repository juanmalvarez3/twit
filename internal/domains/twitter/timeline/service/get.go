package service

import (
	"context"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"go.uber.org/zap"
)

func (s Service) Get(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, bool, error) {
	s.logger.Debug("Obteniendo timeline",
		zap.String("action", actionGet),
		zap.String("user_id", userID),
		zap.Int("limit", limit))

	timeline, cacheHit, err := s.timelineRepo.Get(ctx, userID, limit)
	if err != nil {
		s.logger.Error("Error al obtener timeline",
			zap.String("action", actionGet),
			zap.String("user_id", userID),
			zap.Error(err))
		return dmntimeline.Timeline{}, false, err
	}

	//TODO: Agregar a validaciones
	if len(timeline.Entries) == 0 {
		s.logger.Debug("Timeline vacía",
			zap.String("action", actionGet),
			zap.String("user_id", userID))
		return dmntimeline.Timeline{
			UserID:  userID,
			Entries: []dmntimeline.TimelineEntry{},
		}, false, nil
	}

	s.logger.Debug("Timeline obtenida exitosamente",
		zap.String("action", actionGet),
		zap.String("user_id", userID),
		zap.Int("entries_count", len(timeline.Entries)),
		zap.Bool("cache_hit", cacheHit))

	return timeline, cacheHit, nil
}

func (s Service) GetFromDB(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, error) {
	s.logger.Debug("Obteniendo timeline",
		zap.String("action", actionGet),
		zap.String("user_id", userID),
		zap.Int("limit", limit))

	timeline, err := s.timelineRepo.GetFromDB(ctx, userID, limit)
	if err != nil {
		s.logger.Error("Error al obtener timeline",
			zap.String("action", actionGet),
			zap.String("user_id", userID),
			zap.Error(err))
		return dmntimeline.Timeline{}, err
	}

	if len(timeline.Entries) == 0 {
		s.logger.Debug("Timeline vacía",
			zap.String("action", actionGet),
			zap.String("user_id", userID))
		return dmntimeline.Timeline{
			UserID:  userID,
			Entries: []dmntimeline.TimelineEntry{},
		}, nil
	}

	s.logger.Debug("Timeline obtenida exitosamente",
		zap.String("action", actionGet),
		zap.String("user_id", userID),
		zap.Int("entries_count", len(timeline.Entries)))

	return timeline, nil
}
