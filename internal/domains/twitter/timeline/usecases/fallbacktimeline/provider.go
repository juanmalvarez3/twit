package fallbacktimeline

import (
	srvfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/services"
	srvtweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/services"
	"github.com/juanmalvarez3/twit/pkg/config"
	pkgLogger "github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide(sqsClient SQSClient, cfg *config.Config, logger *pkgLogger.Logger) UseCase {
	publisher := NewTimelinePublisher(
		sqsClient,
		cfg.SQS.RebuildTimelineQueue,
		"rebuild_timeline",
		logger,
	)

	return New(srvfollow.Provide(), srvtweet.Provide(), publisher)
}
