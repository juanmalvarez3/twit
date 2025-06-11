package createtweet_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	createtweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/usecases/createtweet"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/usecases/createtweet/mocks"
)

func TestCreateTweet_Success(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := createtweet.NewUseCase(mockTwtService, mockLogger)

	tweet := &dmntweet.Tweet{UserID: "user-1", Content: "Hello world!"}
	normalized := "Hello world!"
	createdTweet := *tweet
	createdTweet.ID = "twt-123"
	createdTweet.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	mockTwtService.On("Create", mock.Anything, mock.MatchedBy(func(t dmntweet.Tweet) bool {
		return t.UserID == tweet.UserID && t.Content == normalized
	})).Return(createdTweet, nil)

	result, err := uc.CreateTweet(context.Background(), tweet)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdTweet.ID, result.ID)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateTweet_ValidationError(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := createtweet.NewUseCase(mockTwtService, mockLogger)

	tweet := &dmntweet.Tweet{UserID: "user-1", Content: ""}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	result, err := uc.CreateTweet(context.Background(), tweet)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockTwtService.AssertNotCalled(t, "Create")
	mockLogger.AssertExpectations(t)
}

func TestCreateTweet_ServiceError(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := createtweet.NewUseCase(mockTwtService, mockLogger)

	tweet := &dmntweet.Tweet{UserID: "user-1", Content: "Hello world!"}
	serviceErr := errors.New("service error")

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	mockTwtService.On("Create", mock.Anything, mock.MatchedBy(func(t dmntweet.Tweet) bool {
		return t.UserID == tweet.UserID && t.Content == tweet.Content
	})).Return(dmntweet.Tweet{}, serviceErr)

	result, err := uc.CreateTweet(context.Background(), tweet)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceErr, err)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateTweet_ContentTooLong(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := createtweet.NewUseCase(mockTwtService, mockLogger)

	longContent := ""
	for i := 0; i < 281; i++ {
		longContent += "a"
	}
	
	tweet := &dmntweet.Tweet{UserID: "user-1", Content: longContent}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	result, err := uc.CreateTweet(context.Background(), tweet)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "excede el máximo permitido")
	mockTwtService.AssertNotCalled(t, "Create")
	mockLogger.AssertExpectations(t)
}

func TestCreateTweet_ContentNormalization(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := createtweet.NewUseCase(mockTwtService, mockLogger)

	contentWithSpaces := "  Hello world!  "
	normalizedContent := "Hello world!"
	
	tweet := &dmntweet.Tweet{UserID: "user-1", Content: contentWithSpaces}
	createdTweet := dmntweet.Tweet{
		ID:        "twt-123",
		UserID:    "user-1",
		Content:   normalizedContent,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()

	mockTwtService.On("Create", mock.Anything, mock.MatchedBy(func(t dmntweet.Tweet) bool {
		return t.UserID == tweet.UserID && t.Content == normalizedContent
	})).Return(createdTweet, nil)

	result, err := uc.CreateTweet(context.Background(), tweet)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, normalizedContent, result.Content)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateTweet_NoID(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := createtweet.NewUseCase(mockTwtService, mockLogger)

	tweet := &dmntweet.Tweet{UserID: "user-1", Content: "Hello world!"}
	createdTweet := dmntweet.Tweet{
		ID:        "",
		UserID:    "user-1",
		Content:   "Hello world!",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	mockTwtService.On("Create", mock.Anything, mock.Anything).Return(createdTweet, nil)

	result, err := uc.CreateTweet(context.Background(), tweet)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.ID)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateTweet_CanceledContext(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := createtweet.NewUseCase(mockTwtService, mockLogger)

	tweet := &dmntweet.Tweet{UserID: "user-1", Content: "Hello world!"}
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	contextErr := context.Canceled

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	mockTwtService.On("Create", mock.Anything, mock.Anything).Return(dmntweet.Tweet{}, contextErr)

	result, err := uc.CreateTweet(ctx, tweet)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, contextErr, err)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateTweet_EmptyUserID(t *testing.T) {
	mockTwtService := new(mocks.TweetsService)
	mockLogger := new(mocks.Logger)
	uc := createtweet.NewUseCase(mockTwtService, mockLogger)

	tweet := &dmntweet.Tweet{UserID: "", Content: "Hello world!"}
	createdTweet := dmntweet.Tweet{
		ID:        "twt-123",
		UserID:    "",
		Content:   "Hello world!",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()

	mockTwtService.On("Create", mock.Anything, mock.MatchedBy(func(t dmntweet.Tweet) bool {
		return t.UserID == "" && t.Content == "Hello world!"
	})).Return(createdTweet, nil)

	result, err := uc.CreateTweet(context.Background(), tweet)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.UserID)
	mockTwtService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
