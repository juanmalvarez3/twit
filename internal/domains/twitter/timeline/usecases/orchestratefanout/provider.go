package orchestratefanout

import (
	srvfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/services"
	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide(sqsClient SQSClient, cfg *config.Config, log *logger.Logger) UseCase {
	updateTimelinePublisher := UpdateTimelinePublisher(
		sqsClient,
		cfg.SQS.UpdateTimelineQueue,
	)

	followService := srvfollow.Provide()

	useCase := New(
		followService,
		updateTimelinePublisher,
		log,
	)

	return useCase
}
