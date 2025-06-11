package updatetimeline

import "github.com/juanmalvarez3/twit/pkg/logger"

type UseCase struct {
	timelineService TimelineService
	logger          *logger.Logger
}

func New(service TimelineService, log *logger.Logger) *UseCase {
	if log == nil {
		log = logger.DefaultLogger()
	}

	return &UseCase{
		timelineService: service,
		logger:          log,
	}
}
