package service

import (
	"github.com/juanmalvarez3/twit/pkg/logger"
)

const (
	target = "timeline_service"

	actionUpdate = "update"
	actionGet    = "get"
)

type action string

type Service struct {
	timelineRepo Repository
	logger       logger.LoggerInterface
}

func New(timelineRepo Repository, logger logger.LoggerInterface) Service {
	if logger == nil {
		panic("logger cannot be nil")
	}
	return Service{
		timelineRepo: timelineRepo,
		logger:       logger,
	}
}
