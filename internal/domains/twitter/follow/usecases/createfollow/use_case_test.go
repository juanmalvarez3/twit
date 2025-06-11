package createfollow_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/usecases/createfollow"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/usecases/createfollow/mocks"
)

func TestCreateFollow_Success(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := createfollow.NewUseCase(mockService, mockLogger)

	follow := dmnfollow.Follow{
		FollowerID: "user-1",
		FollowedID: "user-2",
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(f dmnfollow.Follow) bool {
		expectedID := "flw-user-1-user-2"
		_, timeErr := time.Parse(time.RFC3339, f.CreatedAt)

		return f.ID == expectedID && f.FollowerID == follow.FollowerID &&
			f.FollowedID == follow.FollowedID && timeErr == nil
	})).Return(nil)

	err := uc.CreateFollow(context.Background(), follow)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateFollow_SameUserError(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := createfollow.NewUseCase(mockService, mockLogger)

	follow := dmnfollow.Follow{
		FollowerID: "user-1",
		FollowedID: "user-1",
	}

	err := uc.CreateFollow(context.Background(), follow)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "no puede seguirse a sí mismo"))
	mockService.AssertNotCalled(t, "Create")
}

func TestCreateFollow_ServiceError(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := createfollow.NewUseCase(mockService, mockLogger)

	follow := dmnfollow.Follow{
		FollowerID: "user-1",
		FollowedID: "user-2",
	}

	serviceErr := errors.New("service error")
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(f dmnfollow.Follow) bool {
		return f.FollowerID == follow.FollowerID && f.FollowedID == follow.FollowedID
	})).Return(serviceErr)

	err := uc.CreateFollow(context.Background(), follow)

	assert.Error(t, err)
	assert.Equal(t, serviceErr, err)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateFollow_WithExistingCreatedAt(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := createfollow.NewUseCase(mockService, mockLogger)

	existingTime := "2025-01-01T12:00:00Z"
	follow := dmnfollow.Follow{
		FollowerID: "user-1",
		FollowedID: "user-2",
		CreatedAt:  existingTime,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(f dmnfollow.Follow) bool {
		return f.CreatedAt == existingTime
	})).Return(nil)

	err := uc.CreateFollow(context.Background(), follow)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateFollow_WithInvalidCreatedAt(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := createfollow.NewUseCase(mockService, mockLogger)

	invalidTime := "2025/01/01 12:00:00"
	follow := dmnfollow.Follow{
		FollowerID: "user-1",
		FollowedID: "user-2",
		CreatedAt:  invalidTime,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(f dmnfollow.Follow) bool {
		return f.CreatedAt == invalidTime
	})).Return(nil)

	err := uc.CreateFollow(context.Background(), follow)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateFollow_CanceledContext(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := createfollow.NewUseCase(mockService, mockLogger)

	follow := dmnfollow.Follow{
		FollowerID: "user-1",
		FollowedID: "user-2",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	contextErr := context.Canceled
	
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(f dmnfollow.Follow) bool {
		return f.FollowerID == follow.FollowerID && f.FollowedID == follow.FollowedID
	})).Return(contextErr)

	err := uc.CreateFollow(ctx, follow)

	assert.Error(t, err)
	assert.Equal(t, contextErr, err)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateFollow_EmptyFollowerID(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := createfollow.NewUseCase(mockService, mockLogger)

	follow := dmnfollow.Follow{
		FollowerID: "",
		FollowedID: "user-2",
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(f dmnfollow.Follow) bool {
		expectedID := "flw--user-2"
		return f.ID == expectedID
	})).Return(nil)

	err := uc.CreateFollow(context.Background(), follow)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateFollow_EmptyFollowedID(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := createfollow.NewUseCase(mockService, mockLogger)

	follow := dmnfollow.Follow{
		FollowerID: "user-1",
		FollowedID: "",
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(f dmnfollow.Follow) bool {
		expectedID := "flw-user-1-"
		return f.ID == expectedID
	})).Return(nil)

	err := uc.CreateFollow(context.Background(), follow)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
