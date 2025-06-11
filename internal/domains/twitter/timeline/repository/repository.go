package repository

import (
	"github.com/juanmalvarez3/twit/pkg/logger"
)

type TimelineRepository struct {
	dynamoDBClient DynamoDBClientInterface
	redisClient    RedisClientInterface
	tableName      string
	logger         *logger.Logger
}

func NewTimelineRepository(
	dynamoDBClient DynamoDBClientInterface,
	redisClient RedisClientInterface,
	tableName string,
	logger *logger.Logger) *TimelineRepository {

	if logger == nil {
		panic("logger cannot be nil")
	}

	namedLogger := logger.Named("timeline_repository")

	return &TimelineRepository{
		dynamoDBClient: dynamoDBClient,
		redisClient:    redisClient,
		tableName:      tableName,
		logger:         namedLogger,
	}
}
