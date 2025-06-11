package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/repository/daos"
	"go.uber.org/zap"
)

func (r *TweetRepository) Create(ctx context.Context, tweet dmntweet.Tweet) error {
	r.logger.Debug("Guardando tweet",
		zap.String("tweet_id", tweet.ID),
		zap.String("user_id", tweet.UserID),
		zap.String("content", tweet.Content),
		zap.String("created_at", tweet.CreatedAt),
		zap.String("table", r.tableName),
	)

	item, err := attributevalue.MarshalMap(daos.ToTweetDAOModel(tweet))
	if err != nil {
		r.logger.Error("Error al serializar tweet para DynamoDB",
			zap.String("tweet_id", tweet.ID),
			zap.String("user_id", tweet.UserID),
			zap.Error(err),
		)
		return err
	}

	_, err = r.dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		r.logger.Error("Error al guardar tweet en DynamoDB",
			zap.String("tweet_id", tweet.ID),
			zap.String("user_id", tweet.UserID),
			zap.String("table", r.tableName),
			zap.Error(err),
		)
		return err
	}

	r.logger.Debug("Tweet guardado exitosamente",
		zap.String("tweet_id", tweet.ID),
		zap.String("user_id", tweet.UserID),
	)
	return nil
}
