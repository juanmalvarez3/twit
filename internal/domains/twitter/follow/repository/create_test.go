package repository

import (
	"context"
	"errors"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/repository/mocks"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepository_Create_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &mocks.MockDBInterface{}
	mockLogger := &mocks.MockLoggerInterface{}
	repo := NewRepository(mockDB, "follows", mockLogger)

	mockDB.On("PutItem", ctx, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	follow := dmnfollow.Follow{
		ID:         "f1",
		FollowerID: "u1",
		FollowedID: "u2",
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	err := repo.Create(ctx, follow)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestRepository_Create_DBError(t *testing.T) {
	ctx := context.Background()
	mockDB := &mocks.MockDBInterface{}
	mockLogger := &mocks.MockLoggerInterface{}
	repo := NewRepository(mockDB, "follows", mockLogger)

	mockDB.On("PutItem", ctx, mock.Anything).Return((*dynamodb.PutItemOutput)(nil), errors.New("db error"))
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	follow := dmnfollow.Follow{
		ID:         "f2",
		FollowerID: "u3",
		FollowedID: "u4",
	}

	err := repo.Create(ctx, follow)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
