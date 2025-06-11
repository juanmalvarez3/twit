package gettimeline

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"go.uber.org/zap"
)

func (u *UseCase) Exec(ctx context.Context, userID string) (dmntimeline.Timeline, error) {
	timeline, cacheHit, err := u.timelineService.Get(ctx, userID, 30)
	if err != nil {
		u.logger.Error("Error obteniendo timeline",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return dmntimeline.Timeline{}, err
	}

	if len(timeline.Entries) == 0 {
		u.logger.Debug("Timeline vacío, solicitando reconstrucción",
			zap.String("user_id", userID),
		)
		err = u.fallbackPublisher.Publish(ctx, userID)
		if err != nil {
			u.logger.Warn("Error publicando solicitud de reconstrucción de timeline",
				zap.String("user_id", userID),
				zap.Error(err),
			)
		}
		return dmntimeline.Timeline{}, dmntimeline.ErrEmptyTimeline
	}

	if !cacheHit {
		u.logger.Debug("Timeline no encontrado en caché, publicando para lazy caching",
			zap.String("user_id", userID),
			zap.Int("entries_count", len(timeline.Entries)),
		)
		err = u.publisher.Publish(ctx, timeline)
		if err != nil {
			u.logger.Warn("Error publicando timeline para lazy caching",
				zap.String("user_id", userID),
				zap.Error(err),
			)
		}
	} else {
		u.logger.Debug("Timeline encontrado en caché",
			zap.String("user_id", userID),
			zap.Int("entries_count", len(timeline.Entries)),
		)
	}

	return timeline, nil
}
