package getfollow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/usecases/getfollow"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/usecases/getfollow/mocks"
)

func TestGetFollowers_Success(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-1"
	expectedFollowers := []string{"user-2", "user-3", "user-4"}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowers", mock.Anything, userID).Return(expectedFollowers, nil)

	followers, err := uc.GetFollowers(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedFollowers, followers)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowers_ServiceError(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-1"
	serviceErr := errors.New("service error")
	var emptySlice []string

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowers", mock.Anything, userID).Return(emptySlice, serviceErr)

	followers, err := uc.GetFollowers(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, serviceErr, err)
	assert.Empty(t, followers)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowing_Success(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-1"
	expectedFollowing := []string{"user-5", "user-6", "user-7"}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowing", mock.Anything, userID).Return(expectedFollowing, nil)

	following, err := uc.GetFollowing(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedFollowing, following)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowing_ServiceError(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-1"
	serviceErr := errors.New("service error")
	var emptySlice []string

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowing", mock.Anything, userID).Return(emptySlice, serviceErr)

	following, err := uc.GetFollowing(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, serviceErr, err)
	assert.Empty(t, following)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowers_EmptyResult(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-with-no-followers"
	emptyFollowers := []string{}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowers", mock.Anything, userID).Return(emptyFollowers, nil)

	followers, err := uc.GetFollowers(context.Background(), userID)

	assert.NoError(t, err)
	assert.Empty(t, followers)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowing_EmptyResult(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-following-nobody"
	emptyFollowing := []string{}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowing", mock.Anything, userID).Return(emptyFollowing, nil)

	following, err := uc.GetFollowing(context.Background(), userID)

	assert.NoError(t, err)
	assert.Empty(t, following)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowers_NilResult(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-with-nil-followers"
	var nilFollowers []string = nil

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowers", mock.Anything, userID).Return(nilFollowers, nil)

	followers, err := uc.GetFollowers(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, followers)
	assert.Empty(t, followers)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowing_NilResult(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-with-nil-following"
	var nilFollowing []string = nil

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowing", mock.Anything, userID).Return(nilFollowing, nil)

	following, err := uc.GetFollowing(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, following)
	assert.Empty(t, following)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowers_LargeResult(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "popular-user"
	largeFollowers := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		largeFollowers[i] = "follower-" + string(rune(i))
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowers", mock.Anything, userID).Return(largeFollowers, nil)

	followers, err := uc.GetFollowers(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, 1000, len(followers))
	assert.Equal(t, largeFollowers, followers)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowing_LargeResult(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-following-many"
	largeFollowing := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		largeFollowing[i] = "followed-" + string(rune(i))
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowing", mock.Anything, userID).Return(largeFollowing, nil)

	following, err := uc.GetFollowing(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, 1000, len(following))
	assert.Equal(t, largeFollowing, following)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowers_CanceledContext(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-1"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	contextErr := context.Canceled
	var emptySlice []string

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowers", mock.Anything, userID).Return(emptySlice, contextErr)

	followers, err := uc.GetFollowers(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, contextErr, err)
	assert.Empty(t, followers)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowing_CanceledContext(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-1"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	contextErr := context.Canceled
	var emptySlice []string

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockService.On("GetFollowing", mock.Anything, userID).Return(emptySlice, contextErr)

	following, err := uc.GetFollowing(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, contextErr, err)
	assert.Empty(t, following)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowers_SpecificLoggerCalls(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-1"
	expectedFollowers := []string{"user-2", "user-3", "user-4"}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	mockService.On("GetFollowers", mock.Anything, userID).Return(expectedFollowers, nil)

	followers, err := uc.GetFollowers(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedFollowers, followers)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetFollowing_SpecificLoggerCalls(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	uc := getfollow.NewUseCase(mockService, mockLogger)

	userID := "user-1"
	expectedFollowing := []string{"user-5", "user-6", "user-7"}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	mockService.On("GetFollowing", mock.Anything, userID).Return(expectedFollowing, nil)

	following, err := uc.GetFollowing(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedFollowing, following)
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
