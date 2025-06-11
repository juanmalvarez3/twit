package gettimeline

import (
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/publisher"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/service"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide(
	cachePublisher publisher.TimelinePublisher,
	rebuildPublisher publisher.RebuildPublisher,
	log logger.LoggerInterface,
) UseCase {
	return New(service.Provide(), cachePublisher, rebuildPublisher, log)
}
