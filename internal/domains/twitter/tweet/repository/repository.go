package repository

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

type TweetRepository struct {
	dynamoDBClient *dynamodb.Client
	tableName      string
	logger         *logger.Logger
}

func NewTweetRepository(
	dynamoDBClient *dynamodb.Client,
	tableName string,
	logger *logger.Logger,
) *TweetRepository {
	return &TweetRepository{
		dynamoDBClient: dynamoDBClient,
		tableName:      tableName,
		logger:         logger,
	}
}

func (r *TweetRepository) serializeLastEvaluatedKey(lastEvaluatedKey map[string]types.AttributeValue) (string, error) {
	bytes, err := json.Marshal(lastEvaluatedKey)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r *TweetRepository) deserializeLastEvaluatedKey(token string) (map[string]types.AttributeValue, error) {
	var lastEvaluatedKey map[string]types.AttributeValue
	err := json.Unmarshal([]byte(token), &lastEvaluatedKey)
	if err != nil {
		return nil, err
	}
	return lastEvaluatedKey, nil
}
