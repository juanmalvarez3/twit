package populatecache

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"go.uber.org/zap"
)

func (uc UseCase) Exec(ctx context.Context, timeline dmntimeline.Timeline) error {
	uc.logger.Info("Iniciando populación de cache",
		zap.String("timeline_id", timeline.UserID),
		zap.Int("entries_count", len(timeline.Entries)),
	)

	timeline, err := uc.timelineService.GetFromDB(ctx, timeline.UserID, 100)
	if err != nil {
		uc.logger.Error("Error al obtener timeline de la base de datos",
			zap.String("timeline_id", timeline.UserID),
			zap.Error(err),
		)
		return err
	}

	err = uc.timelineService.UpdateCache(ctx, timeline)
	if err != nil {
		uc.logger.Error("Error al actualizar cache",
			zap.String("timeline_id", timeline.UserID),
			zap.Error(err),
		)
		return err
	}

	uc.logger.Info("Cache actualizado exitosamente")
	return nil
}
