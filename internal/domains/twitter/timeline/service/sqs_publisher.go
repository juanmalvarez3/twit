package service

import (
	"context"
	"fmt"

	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	evttweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/events"
	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
)

type SQSClient interface {
	Send(ctx context.Context, queueURL string, payload any) error
}

type SQSPublisher struct {
	client    SQSClient
	queueURL  string
	logger    *logger.Logger
	queueName string
}

type RebuildTimelinePayload struct {
	UserID string `json:"user_id"`
}

func (p *SQSPublisher) Publish(ctx context.Context, tweet dmntweet.Tweet, eventType evttweet.EventType) error {
	p.logger.Debug("Preparando envío de mensaje a SQS",
		zap.String("queue_name", p.queueName),
		zap.String("queue_url", p.queueURL),
		zap.String("tweet_id", tweet.ID))

	//TODO: Pasar a constantes
	payload := map[string]interface{}{
		"tweet_id":   tweet.ID,
		"user_id":    tweet.UserID,
		"content":    tweet.Content,
		"created_at": tweet.CreatedAt,
	}

	err := p.client.Send(ctx, p.queueURL, payload)
	if err != nil {
		p.logger.Error("Error al enviar mensaje a SQS",
			zap.String("queue_name", p.queueName),
			zap.String("tweet_id", tweet.ID),
			zap.Error(err))
		return fmt.Errorf("error al enviar mensaje a SQS: %w", err)
	}

	p.logger.Info("Mensaje enviado a SQS exitosamente",
		zap.String("queue_name", p.queueName),
		zap.String("tweet_id", tweet.ID))
	return nil
}
