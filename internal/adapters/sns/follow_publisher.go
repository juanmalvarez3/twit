package sns

import (
	"context"
	"fmt"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain/events"

	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
)

type FollowSNSPublisher struct {
	client   *SNSClient
	topicARN string
	logger   *logger.Logger
}

func NewFollowSNSPublisher(
	client *SNSClient,
	cfg *config.Config,
	logger *logger.Logger,
) *FollowSNSPublisher {
	return &FollowSNSPublisher{
		client:   client,
		topicARN: cfg.SNS.FollowsTopic,
		logger:   logger.With(zap.String("component", "follow_sns_publisher")),
	}
}

func (p *FollowSNSPublisher) Publish(ctx context.Context, event events.Event) error {
	p.logger.Debug("Publicando evento de follow creado en SNS",
		zap.String("follow_id", event.Follow.ID),
		zap.String("follower_id", event.Follow.FollowerID),
		zap.String("followed_id", event.Follow.FollowedID),
		zap.String("topic_arn", p.topicARN))

	payload := events.Event{Follow: event.Follow}

	messageAttributes := map[string]string{
		"event_type":    events.FollowCreatedEventType.String(),
		"resource_type": events.ResourceType,
	}

	err := p.client.PublishMessage(ctx, p.topicARN, payload, messageAttributes)
	if err != nil {
		p.logger.Error("Error publicando evento de follow creado en SNS",
			zap.String("follow_id", event.Follow.ID),
			zap.String("topic_arn", p.topicARN),
			zap.Error(err))
		return fmt.Errorf("error publicando evento de follow: %w", err)
	}

	p.logger.Info("Evento de follow creado publicado exitosamente",
		zap.String("follow_id", event.Follow.ID))
	return nil
}
