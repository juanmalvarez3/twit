package services

import (
	"github.com/juanmalvarez3/twit/pkg/logger"
)

const (
	target = "tweets_service"

	actionCreate = "create"
	actionGet    = "get"
	actionSearch = "search"
)

type action string

type Service struct {
	repository Repository
	publisher  Publisher
	logger     logger.LoggerInterface
}

func NewService(repository Repository, publisher Publisher, log logger.LoggerInterface) Service {
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
