package fallbacktimeline_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/fallbacktimeline"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/fallbacktimeline/mocks"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	dmnoptions "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
)

func TestExec_Success(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockTweetService := new(mocks.TweetService)
	mockPublisher := new(mocks.Publisher)

	uc := fallbacktimeline.New(mockFollowerService, mockTweetService, mockPublisher)

	userID := "user-1"
	following := []string{"user-2", "user-3"}

	tweetsUser2 := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "Hello from user 2",
			CreatedAt: "2025-06-10T22:00:00Z",
		},
	}

	tweetsUser3 := []dmntweet.Tweet{
		{
			ID:        "twt-2",
			UserID:    "user-3",
			Content:   "Hello from user 3",
			CreatedAt: "2025-06-10T23:00:00Z",
		},
	}

	mockFollowerService.On("GetAllFollowing", mock.Anything, userID).Return(following, nil)

	mockTweetService.On("Search", mock.Anything, mock.MatchedBy(func(opts dmnoptions.SearchOptions) bool {
		return opts.Filters.UserID != nil && *opts.Filters.UserID == "user-2"
	})).Return(tweetsUser2, "", nil)

	mockTweetService.On("Search", mock.Anything, mock.MatchedBy(func(opts dmnoptions.SearchOptions) bool {
		return opts.Filters.UserID != nil && *opts.Filters.UserID == "user-3"
	})).Return(tweetsUser3, "", nil)

	mockPublisher.On("Publish", mock.Anything, userID, mock.MatchedBy(func(entries []dmntimeline.TimelineEntry) bool {
		if len(entries) != 2 {
			return false
		}
		if entries[0].AuthorID != "user-3" || entries[0].TweetID != "twt-2" {
			return false
		}
		if entries[1].AuthorID != "user-2" || entries[1].TweetID != "twt-1" {
			return false
		}
		return true
	})).Return(nil)

	err := uc.Exec(context.Background(), userID)

	assert.NoError(t, err)
	mockFollowerService.AssertExpectations(t)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestExec_NoFollowing(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockTweetService := new(mocks.TweetService)
	mockPublisher := new(mocks.Publisher)

	uc := fallbacktimeline.New(mockFollowerService, mockTweetService, mockPublisher)

	userID := "user-1"
	var emptyFollowing []string

	mockFollowerService.On("GetAllFollowing", mock.Anything, userID).Return(emptyFollowing, nil)

	err := uc.Exec(context.Background(), userID)

	assert.NoError(t, err)
	mockFollowerService.AssertExpectations(t)
	mockTweetService.AssertNotCalled(t, "Search")
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_GetFollowingError(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockTweetService := new(mocks.TweetService)
	mockPublisher := new(mocks.Publisher)

	uc := fallbacktimeline.New(mockFollowerService, mockTweetService, mockPublisher)

	userID := "user-1"
	followingErr := errors.New("error getting following")
	var emptyFollowing []string

	mockFollowerService.On("GetAllFollowing", mock.Anything, userID).Return(emptyFollowing, followingErr)

	err := uc.Exec(context.Background(), userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error obteniendo seguidos")
	mockFollowerService.AssertExpectations(t)
	mockTweetService.AssertNotCalled(t, "Search")
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_SearchError(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockTweetService := new(mocks.TweetService)
	mockPublisher := new(mocks.Publisher)

	uc := fallbacktimeline.New(mockFollowerService, mockTweetService, mockPublisher)

	userID := "user-1"
	following := []string{"user-2"}
	searchErr := errors.New("error searching tweets")
	var emptyTweets []dmntweet.Tweet

	mockFollowerService.On("GetAllFollowing", mock.Anything, userID).Return(following, nil)
	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(emptyTweets, "", searchErr)

	err := uc.Exec(context.Background(), userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error obteniendo tweets")
	mockFollowerService.AssertExpectations(t)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_NoTweets(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockTweetService := new(mocks.TweetService)
	mockPublisher := new(mocks.Publisher)

	uc := fallbacktimeline.New(mockFollowerService, mockTweetService, mockPublisher)

	userID := "user-1"
	following := []string{"user-2"}
	var emptyTweets []dmntweet.Tweet

	mockFollowerService.On("GetAllFollowing", mock.Anything, userID).Return(following, nil)
	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(emptyTweets, "", nil)

	err := uc.Exec(context.Background(), userID)

	assert.NoError(t, err)
	mockFollowerService.AssertExpectations(t)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_PublishError(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockTweetService := new(mocks.TweetService)
	mockPublisher := new(mocks.Publisher)

	uc := fallbacktimeline.New(mockFollowerService, mockTweetService, mockPublisher)

	userID := "user-1"
	following := []string{"user-2"}
	tweets := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "Hello from user 2",
			CreatedAt: "2025-06-10T22:00:00Z",
		},
	}
	publishErr := errors.New("error publishing timeline")

	mockFollowerService.On("GetAllFollowing", mock.Anything, userID).Return(following, nil)
	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(tweets, "", nil)
	mockPublisher.On("Publish", mock.Anything, userID, mock.Anything).Return(publishErr)

	err := uc.Exec(context.Background(), userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error publicando timeline reconstruido")
	mockFollowerService.AssertExpectations(t)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestExec_NoPublisher(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockTweetService := new(mocks.TweetService)

	uc := fallbacktimeline.New(mockFollowerService, mockTweetService, nil)

	userID := "user-1"
	following := []string{"user-2"}
	tweets := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "Hello from user 2",
			CreatedAt: "2025-06-10T22:00:00Z",
		},
	}

	mockFollowerService.On("GetAllFollowing", mock.Anything, userID).Return(following, nil)
	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(tweets, "", nil)

	err := uc.Exec(context.Background(), userID)

	assert.NoError(t, err)
	mockFollowerService.AssertExpectations(t)
	mockTweetService.AssertExpectations(t)
}

func TestExec_InvalidTweetDate(t *testing.T) {
	mockFollowerService := new(mocks.FollowerService)
	mockTweetService := new(mocks.TweetService)
	mockPublisher := new(mocks.Publisher)

	uc := fallbacktimeline.New(mockFollowerService, mockTweetService, mockPublisher)

	userID := "user-1"
	following := []string{"user-2"}

	tweets := []dmntweet.Tweet{
		{
			ID:        "twt-1",
			UserID:    "user-2",
			Content:   "Hello from user 2",
			CreatedAt: "invalid-date",
		},
	}

	mockFollowerService.On("GetAllFollowing", mock.Anything, userID).Return(following, nil)
	mockTweetService.On("Search", mock.Anything, mock.Anything).Return(tweets, "", nil)

	mockPublisher.On("Publish", mock.Anything, userID, mock.MatchedBy(func(entries []dmntimeline.TimelineEntry) bool {
		if len(entries) != 1 {
			return false
		}

		now := time.Now().UTC()
		diff := now.Sub(entries[0].CreatedAt)
		return diff < 5*time.Second
	})).Return(nil)

	err := uc.Exec(context.Background(), userID)

	assert.NoError(t, err)
	mockFollowerService.AssertExpectations(t)
	mockTweetService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}
