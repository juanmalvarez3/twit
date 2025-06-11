package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/repository/daos"
	"go.uber.org/zap"
)

func (r *TweetRepository) Get(ctx context.Context, tweetID string) (dmntweet.Tweet, error) {
	r.logger.Debug("Obteniendo tweet por ID",
		zap.String("tweet_id", tweetID),
		zap.String("table_name", r.tableName),
	)

	result, err := r.dynamoDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: tweetID},
		},
	})
	if err != nil {
		r.logger.Error("Error al obtener tweet de DynamoDB",
			zap.String("tweet_id", tweetID),
			zap.Error(err),
		)
		return dmntweet.Tweet{}, err
	}

	if result.Item == nil {
		r.logger.Warn("Tweet no encontrado",
			zap.String("tweet_id", tweetID),
		)
		return dmntweet.Tweet{}, fmt.Errorf("tweet with ID %s not found", tweetID)
	}

	tweet := &daos.TweetDAO{}
	if err := attributevalue.UnmarshalMap(result.Item, tweet); err != nil {
		r.logger.Error("Error al deserializar tweet",
			zap.String("tweet_id", tweetID),
			zap.Error(err),
		)
		return dmntweet.Tweet{}, err
	}

	r.logger.Debug("Tweet obtenido exitosamente",
		zap.String("tweet_id", tweet.ID),
		zap.String("user_id", tweet.UserID),
	)

	return daos.ToTweetModel(*tweet), nil
}
