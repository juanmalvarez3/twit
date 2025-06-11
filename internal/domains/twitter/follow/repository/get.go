package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/repository/daos"
	"go.uber.org/zap"
	"strings"
)

func (r *Repository) Get(ctx context.Context, followID string) (dmnfollow.Follow, error) {
	r.logger.Debug("Obteniendo follow por ID",
		zap.String("follow_id", followID),
		zap.String("table_name", r.tableName),
	)

	parts := strings.Split(followID, "-")
	if len(parts) < 3 {
		err := fmt.Errorf("formato de ID de follow inválido: %s", followID)
		r.logger.Error("Error al parsear ID de follow",
			zap.String("follow_id", followID),
			zap.Error(err),
		)
		return dmnfollow.Follow{}, err
	}

	followerID := parts[1]
	followedID := strings.Join(parts[2:], "-")

	r.logger.Debug("Consultando follow con claves primarias",
		zap.String("follower_id", followerID),
		zap.String("followed_id", followedID),
	)

	result, err := r.dynamoDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"follower_id": &types.AttributeValueMemberS{Value: followerID},
			"followed_id": &types.AttributeValueMemberS{Value: followedID},
		},
	})
	if err != nil {
		r.logger.Error("Error al obtener follow de DynamoDB",
			zap.String("follow_id", followID),
			zap.String("follower_id", followerID),
			zap.String("followed_id", followedID),
			zap.Error(err),
		)
		return dmnfollow.Follow{}, err
	}

	if result.Item == nil {
		r.logger.Warn("Follow no encontrado",
			zap.String("follow_id", followID),
			zap.String("follower_id", followerID),
			zap.String("followed_id", followedID),
		)
		return dmnfollow.Follow{}, fmt.Errorf("follow with ID %s not found", followID)
	}

	follow := &daos.FollowDAO{}
	if err := attributevalue.UnmarshalMap(result.Item, follow); err != nil {
		r.logger.Error("Error al deserializar follow",
			zap.String("follow_id", followID),
			zap.Error(err),
		)
		return dmnfollow.Follow{}, err
	}

	r.logger.Debug("Follow obtenido exitosamente",
		zap.String("follow_id", follow.ID),
		zap.String("follower_id", follow.FollowerID),
		zap.String("followed_id", follow.FollowedID),
	)

	return daos.ToFollowModel(*follow), nil
}

func (r *Repository) GetFollowers(ctx context.Context, followedID string) ([]string, error) {
	r.logger.Debug("Obteniendo seguidores",
		zap.String("followed_id", followedID),
		zap.String("table_name", r.tableName),
	)

	result, err := r.dynamoDBClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("followed_id-index"),
		KeyConditionExpression: aws.String("followed_id = :followedID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":followedID": &types.AttributeValueMemberS{Value: followedID},
		},
		ProjectionExpression: aws.String("follower_id"),
	})
	if err != nil {
		r.logger.Error("Error al consultar seguidores en DynamoDB",
			zap.String("followed_id", followedID),
			zap.Error(err),
		)
		return nil, err
	}

	followers := make([]string, 0)
	for _, item := range result.Items {
		if followerIDAttr, ok := item["follower_id"]; ok {
			if followerIDStr, ok := followerIDAttr.(*types.AttributeValueMemberS); ok {
				followers = append(followers, followerIDStr.Value)
			}
		}
	}

	r.logger.Debug("Seguidores obtenidos exitosamente",
		zap.String("followed_id", followedID),
		zap.Int("count", len(followers)),
	)

	return followers, nil
}

func (r *Repository) GetFollowing(ctx context.Context, followerID string) ([]string, error) {
	r.logger.Debug("Obteniendo usuarios seguidos",
		zap.String("follower_id", followerID),
		zap.String("table_name", r.tableName),
	)

	result, err := r.dynamoDBClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("follower_id-index"),
		KeyConditionExpression: aws.String("follower_id = :followerID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":followerID": &types.AttributeValueMemberS{Value: followerID},
		},
		ProjectionExpression: aws.String("followed_id"),
	})
	if err != nil {
		r.logger.Error("Error al consultar usuarios seguidos en DynamoDB",
			zap.String("follower_id", followerID),
			zap.Error(err),
		)
		return nil, err
	}

	following := make([]string, 0)
	for _, item := range result.Items {
		if followedIDAttr, ok := item["followed_id"]; ok {
			if followedIDStr, ok := followedIDAttr.(*types.AttributeValueMemberS); ok {
				following = append(following, followedIDStr.Value)
			}
		}
	}

	r.logger.Debug("Usuarios seguidos obtenidos exitosamente",
		zap.String("follower_id", followerID),
		zap.Int("count", len(following)),
	)

	return following, nil
}
