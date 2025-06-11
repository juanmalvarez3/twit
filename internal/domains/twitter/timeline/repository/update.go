package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/repository/daos"
	"go.uber.org/zap"
)

func (r *TimelineRepository) Update(ctx context.Context, entry dmntimeline.TimelineEntry, userID string) error {
	da := daos.ToTimelineEntryDAO(userID, entry)
	item, err := attributevalue.MarshalMap(da)
	if err != nil {
		r.logger.Error("Error serializando entrada de timeline", zap.Error(err))
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	}
	_, err = r.dynamoDBClient.PutItem(ctx, input)
	if err != nil {
		r.logger.Error("Error al insertar entrada en DynamoDB", zap.Error(err))
		return err
	}
	return nil
}
