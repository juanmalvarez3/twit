package repository

import (
	"context"
	"time"

	"go.uber.org/zap"
)

func (r *TimelineRepository) SetCache(ctx context.Context, key string, value []byte) error {
	r.logger.Debug("Guardando datos en caché",
		zap.String("key", key))

	ttl := 5 * time.Minute

	err := r.redisClient.Set(ctx, key, value, ttl)
	if err != nil {
		r.logger.Error("Error guardando datos en caché",
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	r.logger.Debug("Datos guardados exitosamente en caché",
		zap.String("key", key),
		zap.Duration("ttl", ttl))

	return nil
}
