package gettweet_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/usecases/gettweet"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/usecases/gettweet/mocks"
)

func TestGetTweet_Success(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := gettweet.NewUseCase(mockTwtService, mockLogger)

	tweetID := "twt-123"
	expectedTweet := dmntweet.Tweet{
		ID:        tweetID,
		UserID:    "user-1",
		Content:   "Hello world!",
		CreatedAt: "2025-06-10T23:00:00Z",
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockTwtService.On("Get", mock.Anything, tweetID).Return(expectedTweet, nil)

	result, err := uc.GetTweet(context.Background(), tweetID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTweet, result)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetTweet_ServiceError(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := gettweet.NewUseCase(mockTwtService, mockLogger)

	tweetID := "twt-123"
	serviceErr := errors.New("service error")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockTwtService.On("Get", mock.Anything, tweetID).Return(dmntweet.Tweet{}, serviceErr)

	result, err := uc.GetTweet(context.Background(), tweetID)

	assert.Error(t, err)
	assert.Equal(t, serviceErr, err)
	assert.Equal(t, dmntweet.Tweet{}, result)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestNewUseCase_NilLogger(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)

	assert.Panics(t, func() {
		gettweet.NewUseCase(mockTwtService, nil)
	})
}

func TestGetTweet_CanceledContext(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := gettweet.NewUseCase(mockTwtService, mockLogger)

	tweetID := "twt-123"

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	contextErr := context.Canceled

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockTwtService.On("Get", mock.Anything, tweetID).Return(dmntweet.Tweet{}, contextErr)

	result, err := uc.GetTweet(ctx, tweetID)

	assert.Error(t, err)
	assert.Equal(t, contextErr, err)
	assert.Equal(t, dmntweet.Tweet{}, result)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetTweet_EmptyID(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := gettweet.NewUseCase(mockTwtService, mockLogger)

	tweetID := ""
	expectedErr := errors.New("tweet ID cannot be empty")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockTwtService.On("Get", mock.Anything, tweetID).Return(dmntweet.Tweet{}, expectedErr)

	result, err := uc.GetTweet(context.Background(), tweetID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, dmntweet.Tweet{}, result)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetTweet_InvalidContent(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := gettweet.NewUseCase(mockTwtService, mockLogger)

	tweetID := "twt-123"

	invalidTweet := dmntweet.Tweet{
		ID:        tweetID,
		UserID:    "user-1",
		Content:   "",
		CreatedAt: "2025-06-10T23:00:00Z",
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockTwtService.On("Get", mock.Anything, tweetID).Return(invalidTweet, nil)

	result, err := uc.GetTweet(context.Background(), tweetID)
	assert.NoError(t, err)
	assert.Equal(t, invalidTweet, result)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetTweet_ContentTooLong(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := gettweet.NewUseCase(mockTwtService, mockLogger)

	tweetID := "twt-123"

	longContent := ""
	for i := 0; i < 281; i++ {
		longContent += "a"
	}

	invalidTweet := dmntweet.Tweet{
		ID:        tweetID,
		UserID:    "user-1",
		Content:   longContent,
		CreatedAt: "2025-06-10T23:00:00Z",
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockTwtService.On("Get", mock.Anything, tweetID).Return(invalidTweet, nil)

	result, err := uc.GetTweet(context.Background(), tweetID)

	assert.NoError(t, err)
	assert.Equal(t, invalidTweet, result)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
