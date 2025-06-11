package processnewfollow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain/events"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/usecases/processnewfollow"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/usecases/processnewfollow/mocks"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
)

func TestProcessNewFollow_Success(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)
	mockPublisher := new(mocks.UpdateTimelineSQSPublisher)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, mockPublisher)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "user-1",
			FollowedID: "user-2",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	tweets := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "Hello world!",
			CreatedAt: "2025-06-10T22:00:00Z",
		},
		{
			ID:        "twt-2",
			UserID:    "user-2",
			Content:   "Another tweet",
			CreatedAt: "2025-06-10T22:30:00Z",
		},
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.MatchedBy(func(opts options.SearchOptions) bool {
		return opts.Filters.UserID != nil && *opts.Filters.UserID == "user-2" &&
			opts.Pagination.Limit == 10 && opts.Pagination.Offset == 0
	})).Return(tweets, "", nil)

	for _, tweet := range tweets {
		mockPublisher.On("Publish", mock.Anything, tweet, "user-2").Return(nil).Once()
	}

	err := uc.ProcessNewFollow(context.Background(), followEvent)

	assert.NoError(t, err)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestProcessNewFollow_SearchError(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)
	mockPublisher := new(mocks.UpdateTimelineSQSPublisher)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, mockPublisher)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "user-1",
			FollowedID: "user-2",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	searchErr := errors.New("error searching tweets")

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.Anything).Return([]dmntweet.Tweet{}, "", searchErr)

	err := uc.ProcessNewFollow(context.Background(), followEvent)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error al buscar tweets")
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestProcessNewFollow_PublishError(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)
	mockPublisher := new(mocks.UpdateTimelineSQSPublisher)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, mockPublisher)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "user-1",
			FollowedID: "user-2",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	tweets := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "Hello world!",
			CreatedAt: "2025-06-10T22:00:00Z",
		},
	}

	publishErr := errors.New("error publishing message")

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(tweets, "", nil)

	mockPublisher.On("Publish", mock.Anything, tweets[0], "user-2").Return(publishErr)

	err := uc.ProcessNewFollow(context.Background(), followEvent)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error al publicar mensaje")
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestProcessNewFollow_NoPublisher(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, nil)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "user-1",
			FollowedID: "user-2",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	tweets := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "Hello world!",
			CreatedAt: "2025-06-10T22:00:00Z",
		},
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(tweets, "", nil)

	err := uc.ProcessNewFollow(context.Background(), followEvent)

	assert.NoError(t, err)
	mockTweetService.AssertExpectations(t)
	mockLogger.AssertCalled(t, "Warn", mock.MatchedBy(func(msg string) bool {
		return msg == "No se ha configurado publicador para actualización de timeline"
	}), mock.Anything)
}

func TestProcessNewFollow_NoTweets(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)
	mockPublisher := new(mocks.UpdateTimelineSQSPublisher)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, mockPublisher)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "user-1",
			FollowedID: "user-2",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	emptyTweets := []dmntweet.Tweet{}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(emptyTweets, "", nil)

	err := uc.ProcessNewFollow(context.Background(), followEvent)

	assert.NoError(t, err)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestProcessNewFollow_CanceledContext(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)
	mockPublisher := new(mocks.UpdateTimelineSQSPublisher)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, mockPublisher)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "user-1",
			FollowedID: "user-2",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	contextErr := context.Canceled
	
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.Anything).Return([]dmntweet.Tweet{}, "", contextErr)

	err := uc.ProcessNewFollow(ctx, followEvent)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error al buscar tweets")
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestProcessNewFollow_MultiplePublishErrors(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)
	mockPublisher := new(mocks.UpdateTimelineSQSPublisher)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, mockPublisher)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "user-1",
			FollowedID: "user-2",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	tweets := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "First tweet",
			CreatedAt: "2025-06-10T22:00:00Z",
		},
		{
			ID:        "twt-2",
			UserID:    "user-2",
			Content:   "Second tweet",
			CreatedAt: "2025-06-10T22:30:00Z",
		},
	}

	publishErr := errors.New("error publishing message")

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(tweets, "", nil)

	mockPublisher.On("Publish", mock.Anything, tweets[0], "user-2").Return(publishErr)

	err := uc.ProcessNewFollow(context.Background(), followEvent)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error al publicar mensaje")
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
	mockPublisher.AssertNumberOfCalls(t, "Publish", 1)
}

func TestProcessNewFollow_EmptyFollowerID(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)
	mockPublisher := new(mocks.UpdateTimelineSQSPublisher)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, mockPublisher)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "",
			FollowedID: "user-2",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	tweets := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "Hello world!",
			CreatedAt: "2025-06-10T22:00:00Z",
		},
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.MatchedBy(func(opts options.SearchOptions) bool {
		return opts.Filters.UserID != nil && *opts.Filters.UserID == "user-2"
	})).Return(tweets, "", nil)

	mockPublisher.On("Publish", mock.Anything, tweets[0], "user-2").Return(nil)

	err := uc.ProcessNewFollow(context.Background(), followEvent)

	assert.NoError(t, err)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestProcessNewFollow_EmptyFollowedID(t *testing.T) {
	mockService := new(mocks.Service)
	mockLogger := new(mocks.Logger)
	mockTweetService := new(mocks.TweetsService)
	mockPublisher := new(mocks.UpdateTimelineSQSPublisher)

	uc := processnewfollow.NewUseCase(mockService, mockLogger, mockTweetService, mockPublisher)

	followEvent := events.FollowCreatedEvent{
		Type: "FollowCreatedEventType",
		Follow: events.Follow{
			ID:         "flw-123",
			FollowerID: "user-1",
			FollowedID: "",
			CreatedAt:  "2025-06-10T23:00:00Z",
		},
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()

	mockTweetService.On("Search", mock.Anything, mock.MatchedBy(func(opts options.SearchOptions) bool {
		return opts.Filters.UserID != nil && *opts.Filters.UserID == ""
	})).Return([]dmntweet.Tweet{}, "", nil)

	err := uc.ProcessNewFollow(context.Background(), followEvent)

	assert.NoError(t, err)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}
