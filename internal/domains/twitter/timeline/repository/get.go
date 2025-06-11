package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"go.uber.org/zap"
)

func (r *TimelineRepository) Get(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, bool, error) {
	//TODO: Pasar a constante o generate key
	cacheKey := "timeline:" + userID
	cacheData, err := r.redisClient.Get(ctx, cacheKey)

	//TODO: Refactor maybe
	if err == nil && len(cacheData) > 0 {
		var timeline dmntimeline.Timeline
		err = json.Unmarshal(cacheData, &timeline)
		if err == nil && len(timeline.Entries) > 0 {
			r.logger.Debug("Timeline obtenida desde caché",
				zap.String("user_id", userID),
				zap.Int("entries_count", len(timeline.Entries)))
			return timeline, true, nil
		}

		if err != nil {
			r.logger.Warn("Error deserializando timeline desde caché",
				zap.String("user_id", userID),
				zap.Error(err))
		}
	}
	r.logger.Debug("Consultando timeline en DynamoDB",
		zap.String("user_id", userID),
		zap.String("table_name", r.tableName))

	keyEx := expression.Key("user_id").Equal(expression.Value(userID))

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		r.logger.Error("Error construyendo expresión para DynamoDB",
			zap.String("user_id", userID),
			zap.Error(err))
		return dmntimeline.Timeline{}, false, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit)),
		ScanIndexForward:          aws.Bool(false),
	}

	result, err := r.dynamoDBClient.Query(ctx, queryInput)
	if err != nil {
		r.logger.Error("Error consultando timeline en DynamoDB",
			zap.String("user_id", userID),
			zap.String("table_name", r.tableName),
			zap.Error(err))
		return dmntimeline.Timeline{}, false, err
	}

	if len(result.Items) == 0 {
		r.logger.Debug("No se encontraron entradas para la timeline",
			zap.String("user_id", userID))
		return dmntimeline.Timeline{
			UserID:  userID,
			Entries: []dmntimeline.TimelineEntry{},
		}, false, nil
	}

	entries := make([]dmntimeline.TimelineEntry, 0, len(result.Items))
	for _, item := range result.Items {
		var entry dmntimeline.TimelineEntry

		tweetID, _ := item["tweet_id"].(*types.AttributeValueMemberS)
		authorID, _ := item["author_id"].(*types.AttributeValueMemberS)
		content, _ := item["content"].(*types.AttributeValueMemberS)
		createdAtStr, _ := item["created_at"].(*types.AttributeValueMemberS)

		createdAt := time.Now()
		if createdAtStr != nil {
			parsedTime, err := time.Parse(time.RFC3339, createdAtStr.Value)
			if err == nil {
				createdAt = parsedTime
			}
		}

		if tweetID != nil && authorID != nil {
			entry = dmntimeline.TimelineEntry{
				TweetID:   tweetID.Value,
				AuthorID:  authorID.Value,
				Content:   content.Value,
				CreatedAt: createdAt,
			}
			entries = append(entries, entry)
		}
	}

	timeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: entries,
	}

	r.logger.Debug("Timeline obtenida exitosamente desde DynamoDB",
		zap.String("user_id", userID),
		zap.Int("entries_count", len(entries)))

	return timeline, false, nil
}

func (r *TimelineRepository) GetFromDB(ctx context.Context, userID string, limit int) (dmntimeline.Timeline, error) {
	r.logger.Debug("Consultando timeline en DynamoDB",
		zap.String("user_id", userID),
		zap.String("table_name", r.tableName))

	keyEx := expression.Key("user_id").Equal(expression.Value(userID))

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		r.logger.Error("Error construyendo expresión para DynamoDB",
			zap.String("user_id", userID),
			zap.Error(err))
		return dmntimeline.Timeline{}, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit)),
		ScanIndexForward:          aws.Bool(false),
	}

	result, err := r.dynamoDBClient.Query(ctx, queryInput)
	if err != nil {
		r.logger.Error("Error consultando timeline en DynamoDB",
			zap.String("user_id", userID),
			zap.String("table_name", r.tableName),
			zap.Error(err))
		return dmntimeline.Timeline{}, err
	}

	if len(result.Items) == 0 {
		r.logger.Debug("No se encontraron entradas para la timeline",
			zap.String("user_id", userID))
		return dmntimeline.Timeline{
			UserID:  userID,
			Entries: []dmntimeline.TimelineEntry{},
		}, nil
	}

	entries := make([]dmntimeline.TimelineEntry, 0, len(result.Items))
	for _, item := range result.Items {
		var entry dmntimeline.TimelineEntry

		tweetID, _ := item["tweet_id"].(*types.AttributeValueMemberS)
		authorID, _ := item["author_id"].(*types.AttributeValueMemberS)
		content, _ := item["content"].(*types.AttributeValueMemberS)
		createdAtStr, _ := item["created_at"].(*types.AttributeValueMemberS)

		createdAt := time.Now()
		if createdAtStr != nil {
			parsedTime, err := time.Parse(time.RFC3339, createdAtStr.Value)
			if err == nil {
				createdAt = parsedTime
			}
		}

		if tweetID != nil && authorID != nil {
			entry = dmntimeline.TimelineEntry{
				TweetID:   tweetID.Value,
				AuthorID:  authorID.Value,
				Content:   content.Value,
				CreatedAt: createdAt,
			}
			entries = append(entries, entry)
		}
	}

	timeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: entries,
	}

	r.logger.Debug("Timeline obtenida exitosamente desde DynamoDB",
		zap.String("user_id", userID),
		zap.Int("entries_count", len(entries)))

	return timeline, nil
}
