package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"go.uber.org/zap"
	"time"
)

type TweetDAO struct {
	ID        string    `dynamodbav:"id"`
	UserID    string    `dynamodbav:"user_id"`
	Content   string    `dynamodbav:"content"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}

func (r *TweetRepository) Search(ctx context.Context, userID string, limit int, lastEvaluatedKey string) ([]dmntweet.Tweet, string, error) {
	r.logger.Debug("Buscando tweets por usuario",
		zap.String("user_id", userID),
		zap.Int("limit", limit),
		zap.String("table_name", r.tableName),
		zap.String("index_name", "user_id-created_at-index"),
		zap.Bool("has_pagination_token", lastEvaluatedKey != ""),
	)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("user_id-created_at-index"),
		KeyConditionExpression: aws.String("user_id = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{Value: userID},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(limit)),
	}

	if lastEvaluatedKey != "" {
		exclusiveStartKey, err := r.deserializeLastEvaluatedKey(lastEvaluatedKey)
		if err != nil {
			r.logger.Error("Error al deserializar token de paginación",
				zap.String("token", lastEvaluatedKey),
				zap.Error(err),
			)
			return nil, "", err
		}
		input.ExclusiveStartKey = exclusiveStartKey
	}

	result, err := r.dynamoDBClient.Query(ctx, input)
	if err != nil {
		r.logger.Error("Error al consultar tweets en DynamoDB",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, "", err
	}

	r.logger.Debug("Consulta de tweets completada",
		zap.String("user_id", userID),
		zap.Int("items_count", len(result.Items)),
		zap.Bool("has_more", result.LastEvaluatedKey != nil),
	)

	var daos []TweetDAO
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &daos); err != nil {
		r.logger.Error("Error al deserializar tweets de DynamoDB",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, "", err
	}

	tweets := make([]dmntweet.Tweet, len(daos))
	for i, dao := range daos {
		tweets[i] = dmntweet.Tweet{
			ID:        dao.ID,
			UserID:    dao.UserID,
			Content:   dao.Content,
			CreatedAt: dao.CreatedAt.Format(time.RFC3339),
		}
	}

	var nextToken string
	if result.LastEvaluatedKey != nil {
		nextToken, err = r.serializeLastEvaluatedKey(result.LastEvaluatedKey)
		if err != nil {
			r.logger.Error("Error al serializar token de paginación",
				zap.Error(err),
			)
			return nil, "", err
		}
	}

	r.logger.Debug("Tweets encontrados exitosamente",
		zap.String("user_id", userID),
		zap.Int("count", len(tweets)),
		zap.Bool("has_next_page", nextToken != ""),
	)

	return tweets, nextToken, nil
}
