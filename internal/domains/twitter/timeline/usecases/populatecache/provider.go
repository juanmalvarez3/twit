package populatecache

import (
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/service"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide(
	log logger.LoggerInterface,
) UseCase {
	return New(service.Provide(), log)
}
