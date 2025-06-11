package populatecache

const componentName = "populatecache_usecase"

type UseCase struct {
	timelineService TimelineService
	logger          Logger
}

func New(
	timelineService TimelineService,
	logger Logger,
) UseCase {
	if logger == nil {
		panic("logger cannot be nil")
	}

	return UseCase{
		timelineService: timelineService,
		logger:          logger,
	}
}
