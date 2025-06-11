package processnewfollow

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"

	srvfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/services"
	srvtweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/services"

	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

type SQSClientAdapter interface {
	Send(ctx context.Context, queueURL string, payload any) error
}

type UpdateTimelinePublisher struct {
	client   SQSClientAdapter
	queueURL string
	logger   *logger.Logger
}

func (p *UpdateTimelinePublisher) Publish(ctx context.Context, tweet dmntweet.Tweet, userID string) error {
	payload := struct {
		Tweet      dmntweet.Tweet `json:"tweet"`
		FollowerID string         `json:"user_id"`
	}{
		FollowerID: userID,
		Tweet:      tweet,
	}

	return p.client.Send(ctx, p.queueURL, payload)
}

func NewUpdateTimelinePublisher(client SQSClientAdapter, queueURL string, logger *logger.Logger) *UpdateTimelinePublisher {
	return &UpdateTimelinePublisher{
		client:   client,
		queueURL: queueURL,
		logger:   logger,
	}
}

func Provide(
	sqsClient SQSClientAdapter,
	cfg *config.Config,
	log *logger.Logger,
) *UseCase {
	updateTimelinePublisher := NewUpdateTimelinePublisher(
		sqsClient,
		cfg.SQS.UpdateTimelineQueue,
		log,
	)

	followService := srvfollow.Provide()
	tweetService := srvtweet.Provide()

	useCase := NewUseCase(
		followService,
		log,
		tweetService,
		updateTimelinePublisher,
	)

	return &useCase
}
