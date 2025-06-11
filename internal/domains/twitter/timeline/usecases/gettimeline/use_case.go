package gettimeline

const componentName = "gettimeline_usecase"

type UseCase struct {
	timelineService   TimelineService
	publisher         Publisher
	fallbackPublisher FallbackRebuildTimelinePublisherService
	logger            Logger
}

func New(
	service TimelineService,
	publisher Publisher,
	fallbackPublisher FallbackRebuildTimelinePublisherService,
	logger Logger,
) UseCase {
	if logger == nil {
		panic("logger cannot be nil")
	}

	return UseCase{
		timelineService:   service,
		publisher:         publisher,
		fallbackPublisher: fallbackPublisher,
		logger:            logger,
	}
}
