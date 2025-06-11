package queue

import (
	"context"
	"encoding/json"
	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
)

type RebuildTimelinePublisher struct {
	adapter  *Adapter
	queueURL string
	logger   logger.LoggerInterface
}

func NewRebuildTimelinePublisher(adapter *Adapter, queueURL string, logger logger.LoggerInterface) *RebuildTimelinePublisher {
	return &RebuildTimelinePublisher{
		adapter:  adapter,
		queueURL: queueURL,
		logger:   logger,
	}
}

func (p *RebuildTimelinePublisher) Publish(ctx context.Context, userID string) error {
	payload := map[string]string{"user_id": userID}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("Error serializando solicitud de reconstrucción", zap.Error(err))
		return err
	}

	err = p.adapter.Send(ctx, p.queueURL, string(jsonPayload))
	if err != nil {
		p.logger.Error("Error publicando solicitud de reconstrucción", zap.Error(err))
		return err
	}

	return nil
}
