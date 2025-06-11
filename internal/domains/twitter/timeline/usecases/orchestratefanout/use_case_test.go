package orchestratefanout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/orchestratefanout"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/orchestratefanout/mocks"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
)

func TestExec_Success(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockPublisher := new(mocks.Publisher)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	uc := orchestratefanout.New(mockFollowerService, mockPublisher, mockLogger)

	now := time.Now().UTC()
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "Hello world!",
		CreatedAt: now.Format(time.RFC3339),
	}
	followers := []string{"follower-1", "follower-2", "follower-3"}

	mockFollowerService.On("GetFollowers", mock.Anything, tweet.UserID).Return(followers, nil)

	for _, followerID := range followers {
		mockPublisher.On("Publish", mock.Anything, mock.MatchedBy(func(t dmntweet.Tweet) bool {
			return t.ID == tweet.ID && t.UserID == tweet.UserID && t.Content == tweet.Content
		}), followerID).Return(nil)
	}

	err := uc.Exec(context.Background(), tweet)

	assert.NoError(t, err)
	mockFollowerService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestExec_NoFollowers(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockPublisher := new(mocks.Publisher)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	uc := orchestratefanout.New(mockFollowerService, mockPublisher, mockLogger)

	now := time.Now().UTC()
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "Hello world!",
		CreatedAt: now.Format(time.RFC3339),
	}
	emptyFollowers := []string{}

	mockFollowerService.On("GetFollowers", mock.Anything, tweet.UserID).Return(emptyFollowers, nil)

	err := uc.Exec(context.Background(), tweet)

	assert.NoError(t, err)
	mockFollowerService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_FollowerServiceError(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockPublisher := new(mocks.Publisher)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	uc := orchestratefanout.New(mockFollowerService, mockPublisher, mockLogger)

	now := time.Now().UTC()
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "Hello world!",
		CreatedAt: now.Format(time.RFC3339),
	}
	serviceErr := errors.New("error obteniendo seguidores")

	mockFollowerService.On("GetFollowers", mock.Anything, tweet.UserID).Return([]string{}, serviceErr)

	err := uc.Exec(context.Background(), tweet)

	assert.Error(t, err)
	assert.Equal(t, serviceErr, err)
	mockFollowerService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_PublisherError(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockPublisher := new(mocks.Publisher)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	uc := orchestratefanout.New(mockFollowerService, mockPublisher, mockLogger)

	now := time.Now().UTC()
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "Hello world!",
		CreatedAt: now.Format(time.RFC3339),
	}
	followers := []string{"follower-1", "follower-2", "follower-3"}
	publisherErr := errors.New("error publicando evento")

	mockFollowerService.On("GetFollowers", mock.Anything, tweet.UserID).Return(followers, nil)

	mockPublisher.On("Publish", mock.Anything, mock.MatchedBy(func(t dmntweet.Tweet) bool {
		return t.ID == tweet.ID && t.UserID == tweet.UserID && t.Content == tweet.Content
	}), "follower-1").Return(publisherErr)

	mockPublisher.On("Publish", mock.Anything, mock.MatchedBy(func(t dmntweet.Tweet) bool {
		return t.ID == tweet.ID && t.UserID == tweet.UserID && t.Content == tweet.Content
	}), "follower-2").Return(nil)

	mockPublisher.On("Publish", mock.Anything, mock.MatchedBy(func(t dmntweet.Tweet) bool {
		return t.ID == tweet.ID && t.UserID == tweet.UserID && t.Content == tweet.Content
	}), "follower-3").Return(nil)

	err := uc.Exec(context.Background(), tweet)

	assert.NoError(t, err)
	mockFollowerService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestNew_WithNilLogger(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockPublisher := new(mocks.Publisher)

	assert.Panics(t, func() {
		orchestratefanout.New(mockFollowerService, mockPublisher, nil)
	}, "Se espera un pánico cuando el logger es nil")
}
