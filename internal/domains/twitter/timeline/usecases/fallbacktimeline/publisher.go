package fallbacktimeline

import (
	"context"
	"fmt"
	"github.com/juanmalvarez3/twit/pkg/logger"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"go.uber.org/zap"
)

type SQSClient interface {
	Send(ctx context.Context, queueURL string, payload any) error
}

type TimelinePublisher struct {
	client    SQSClient
	queueURL  string
	logger    *logger.Logger
	queueName string
}

func NewTimelinePublisher(
	client SQSClient,
	queueURL string,
	queueName string,
	logger *logger.Logger,
) *TimelinePublisher {
	return &TimelinePublisher{
		client:    client,
		queueURL:  queueURL,
		queueName: queueName,
		logger:    logger,
	}
}

type TimelinePublishPayload struct {
	UserID  string                      `json:"user_id"`
	Entries []dmntimeline.TimelineEntry `json:"entries"`
}

func (p *TimelinePublisher) Publish(ctx context.Context, userID string, entries []dmntimeline.TimelineEntry) error {
	p.logger.Debug("Preparando envío de mensaje a SQS",
		zap.String("queue_name", p.queueName),
		zap.String("user_id", userID),
		zap.Int("entries_count", len(entries)))

	payload := TimelinePublishPayload{
		UserID:  userID,
		Entries: entries,
	}

	err := p.client.Send(ctx, p.queueURL, payload)
	if err != nil {
		p.logger.Error("Error al enviar mensaje a SQS",
			zap.String("queue_name", p.queueName),
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("error al enviar mensaje a SQS: %w", err)
	}

	p.logger.Info("Mensaje enviado a SQS exitosamente",
		zap.String("queue_name", p.queueName),
		zap.String("user_id", userID))
	return nil
}
