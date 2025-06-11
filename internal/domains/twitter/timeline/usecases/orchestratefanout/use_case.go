package orchestratefanout

const componentName = "orchestratefanout_usecase"

type UseCase struct {
	followerService FollowerService
	publisher       Publisher
	logger          Logger
}

func New(
	followerService FollowerService,
	publisher Publisher,
	logger Logger,
) UseCase {
	if logger == nil {
		panic("logger cannot be nil")
	}

	return UseCase{
		followerService: followerService,
		publisher:       publisher,
		logger:          logger,
	}
}
