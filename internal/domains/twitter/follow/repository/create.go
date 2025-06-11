package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/repository/daos"
	"go.uber.org/zap"
	"time"
)

func (r *Repository) Create(ctx context.Context, follow dmnfollow.Follow) error {

	follow.ID = dbPrefix + follow.FollowerID + "-" + follow.FollowedID
	r.logger.Debug("Guardando follow",
		zap.String("follow_id", follow.ID),
		zap.String("follower_id", follow.FollowerID),
		zap.String("followed_id", follow.FollowedID),
		zap.String("table", r.tableName),
	)

	if follow.CreatedAt == "" {
		follow.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	item, err := attributevalue.MarshalMap(daos.ToFollowDAOModel(follow))
	if err != nil {
		r.logger.Error("Error al serializar follow para DynamoDB",
			zap.String("follow_id", follow.ID),
			zap.String("follower_id", follow.FollowerID),
			zap.String("followed_id", follow.FollowedID),
			zap.Error(err),
		)
		return err
	}

	_, err = r.dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		r.logger.Error("Error al guardar follow en DynamoDB",
			zap.String("follow_id", follow.ID),
			zap.String("follower_id", follow.FollowerID),
			zap.String("followed_id", follow.FollowedID),
			zap.String("table", r.tableName),
			zap.Error(err),
		)
		return err
	}

	r.logger.Debug("Follow guardado exitosamente",
		zap.String("follow_id", follow.ID),
		zap.String("follower_id", follow.FollowerID),
		zap.String("followed_id", follow.FollowedID),
	)
	return nil
}
