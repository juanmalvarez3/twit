package services

import "github.com/juanmalvarez3/twit/pkg/logger"

const target = "follow_service"

type Service struct {
	repository Repository
	publisher  Publisher
	logger     *logger.Logger
}

func New(repository Repository, publisher Publisher, log logger.LoggerInterface) Service {
	if log == nil {
		panic("logger cannot be nil")
	}

	serviceLogger := log.Named(target)

	return Service{
		repository: repository,
		publisher:  publisher,
		logger:     serviceLogger,
	}
}
