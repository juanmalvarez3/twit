package queue

import (
	"context"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
)

type PopulateTimelineCachePublisher struct {
	adapter  *Adapter
	queueURL string
	logger   logger.LoggerInterface
}

func NewPopulateTimelineCachePublisher(adapter *Adapter, queueURL string, logger logger.LoggerInterface) *PopulateTimelineCachePublisher {
	return &PopulateTimelineCachePublisher{
		adapter:  adapter,
		queueURL: queueURL,
		logger:   logger,
	}
}

func (p *PopulateTimelineCachePublisher) Publish(ctx context.Context, timeline dmntimeline.Timeline) error {
	err := p.adapter.Send(ctx, p.queueURL, timeline)
	if err != nil {
		p.logger.Error("Error publicando timeline en cola", zap.Error(err))
		return err
	}

	return nil
}
