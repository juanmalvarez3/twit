package service

import (
	"context"
	"fmt"

	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
)

type RebuildTimelinePublisher struct {
	client    SQSClient
	queueURL  string
	logger    *logger.Logger
	queueName string
}

func NewRebuildTimelinePublisher(
	client SQSClient,
	queueURL string,
	queueName string,
	logger *logger.Logger,
) *RebuildTimelinePublisher {
	return &RebuildTimelinePublisher{
		client:    client,
		queueURL:  queueURL,
		queueName: queueName,
		logger:    logger,
	}
}

func (p *RebuildTimelinePublisher) PublishRebuildRequest(ctx context.Context, userID string) error {
	p.logger.Debug("Solicitando reconstrucción de timeline",
		zap.String("queue_name", p.queueName),
		zap.String("queue_url", p.queueURL),
		zap.String("user_id", userID))

	payload := RebuildTimelinePayload{
		UserID: userID,
	}

	err := p.client.Send(ctx, p.queueURL, payload)
	if err != nil {
		p.logger.Error("Error al solicitar reconstrucción de timeline",
			zap.String("queue_name", p.queueName),
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("error al solicitar reconstrucción de timeline: %w", err)
	}

	p.logger.Info("Solicitud de reconstrucción de timeline enviada",
		zap.String("queue_name", p.queueName),
		zap.String("user_id", userID))
	return nil
}
