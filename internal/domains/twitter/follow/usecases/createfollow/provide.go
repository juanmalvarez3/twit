package createfollow

import (
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/services"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide(log *logger.Logger) UseCase {
	return NewUseCase(services.Provide(), log)
}
