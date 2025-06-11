package createfollow

type UseCase struct {
	service Service
	logger  Logger
}

func NewUseCase(
	service Service,
	logger Logger,
) UseCase {
	return UseCase{
		service: service,
		logger:  logger,
	}
}
