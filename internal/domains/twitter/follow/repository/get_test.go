package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepository_Get_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &mocks.MockDBInterface{}
	mockLogger := &mocks.MockLoggerInterface{}
	repo := NewRepository(mockDB, "follows", mockLogger)

	item := map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: "f1"},
		"follower_id": &types.AttributeValueMemberS{Value: "u1"},
		"followed_id": &types.AttributeValueMemberS{Value: "u2"},
		"created_at":  &types.AttributeValueMemberS{Value: "2023-01-01T00:00:00Z"},
	}
	mockDB.On("GetItem", ctx, mock.Anything).Return(&dynamodb.GetItemOutput{Item: item}, nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	follow, err := repo.Get(ctx, "f1")
	assert.NoError(t, err)
	assert.Equal(t, "f1", follow.ID)
	assert.Equal(t, "u1", follow.FollowerID)
	assert.Equal(t, "u2", follow.FollowedID)
	mockDB.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRepository_Get_NotFound(t *testing.T) {
	ctx := context.Background()
	mockDB := &mocks.MockDBInterface{}
	mockLogger := &mocks.MockLoggerInterface{}
	repo := NewRepository(mockDB, "follows", mockLogger)

	mockDB.On("GetItem", ctx, mock.Anything).Return(&dynamodb.GetItemOutput{Item: nil}, nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	_, err := repo.Get(ctx, "notfound")
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRepository_Get_DBError(t *testing.T) {
	ctx := context.Background()
	mockDB := &mocks.MockDBInterface{}
	mockLogger := &mocks.MockLoggerInterface{}
	repo := NewRepository(mockDB, "follows", mockLogger)

	mockDB.On("GetItem", ctx, mock.Anything).Return((*dynamodb.GetItemOutput)(nil), errors.New("db error"))
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	_, err := repo.Get(ctx, "f1")
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRepository_Get_UnmarshalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &mocks.MockDBInterface{}
	mockLogger := &mocks.MockLoggerInterface{}
	repo := NewRepository(mockDB, "follows", mockLogger)

	item := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberN{Value: "not_a_string"},
	}
	mockDB.On("GetItem", ctx, mock.Anything).Return(&dynamodb.GetItemOutput{Item: item}, nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	_, err := repo.Get(ctx, "f1")
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
