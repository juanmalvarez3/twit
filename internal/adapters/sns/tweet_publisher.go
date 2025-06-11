package sns

import (
	"context"
	"fmt"

	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/events"
	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
)

type TweetSNSPublisher struct {
	client   *SNSClient
	topicARN string
	logger   *logger.Logger
}

func NewTweetSNSPublisher(
	client *SNSClient,
	cfg *config.Config,
	logger *logger.Logger,
) *TweetSNSPublisher {
	return &TweetSNSPublisher{
		client:   client,
		topicARN: cfg.SNS.TweetsTopic,
		logger:   logger.With(zap.String("component", "tweet_sns_publisher")),
	}
}

func (p *TweetSNSPublisher) Publish(ctx context.Context, event events.Event) error {
	p.logger.Debug("Publicando evento de tweet en SNS",
		zap.String("tweet_id", event.Tweet.ID),
		zap.String("event_type", event.Type.String()),
		zap.String("topic_arn", p.topicARN))

	var payload interface{}

	switch event.Type {
	case events.TweetCreatedEventType:
		payload = events.NewTweetCreatedEvent(event.Tweet)
	default:
		p.logger.Error("Tipo de evento desconocido",
			zap.String("tweet_id", event.Tweet.ID),
			zap.String("event_type", event.Type.String()))
		return fmt.Errorf("tipo de evento desconocido: %s", event.Type.String())
	}

	messageAttributes := map[string]string{
		"event_type":    event.Type.String(),
		"resource_type": events.ResourceType,
	}

	err := p.client.PublishMessage(ctx, p.topicARN, payload, messageAttributes)
	if err != nil {
		p.logger.Error("Error publicando evento de tweet en SNS",
			zap.String("tweet_id", event.Tweet.ID),
			zap.String("topic_arn", p.topicARN),
			zap.Error(err))
		return fmt.Errorf("error publicando evento de tweet: %w", err)
	}

	p.logger.Info("Evento de tweet publicado exitosamente",
		zap.String("tweet_id", event.Tweet.ID),
		zap.String("event_type", event.Type.String()))
	return nil
}
