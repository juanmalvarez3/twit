package updatetimeline_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/updatetimeline"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/updatetimeline/mocks"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

func TestExec_Success(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	testLogger, _ := logger.New("debug", "test")

	uc := updatetimeline.New(mockTimelineService, testLogger)

	userID := "user-1"
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "Hello world!",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	mockTimelineService.On("Update", mock.Anything, mock.MatchedBy(func(entry dmntimeline.TimelineEntry) bool {
		return entry.TweetID == tweet.ID &&
			entry.AuthorID == tweet.UserID &&
			entry.Content == tweet.Content
	}), userID).Return(nil)

	err := uc.Exec(context.Background(), tweet, userID)

	assert.NoError(t, err)
	mockTimelineService.AssertExpectations(t)
}

func TestExec_ServiceError(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	testLogger, _ := logger.New("debug", "test")

	uc := updatetimeline.New(mockTimelineService, testLogger)

	userID := "user-1"
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "Hello world!",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	expectedErr := errors.New("error al actualizar timeline")

	mockTimelineService.On("Update", mock.Anything, mock.Anything, userID).Return(expectedErr)

	err := uc.Exec(context.Background(), tweet, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockTimelineService.AssertExpectations(t)
}

func TestExec_EmptyContent(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	testLogger, _ := logger.New("debug", "test")

	uc := updatetimeline.New(mockTimelineService, testLogger)

	userID := "user-1"
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	mockTimelineService.On("Update", mock.Anything, mock.MatchedBy(func(entry dmntimeline.TimelineEntry) bool {
		return entry.TweetID == tweet.ID &&
			entry.AuthorID == tweet.UserID &&
			entry.Content == ""
	}), userID).Return(nil)

	err := uc.Exec(context.Background(), tweet, userID)

	assert.NoError(t, err)
	mockTimelineService.AssertExpectations(t)
}

func TestExec_InvalidDate(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	testLogger, _ := logger.New("debug", "test")

	uc := updatetimeline.New(mockTimelineService, testLogger)

	userID := "user-1"
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "Hello world!",
		CreatedAt: "fecha-invalida",
	}

	mockTimelineService.On("Update", mock.Anything, mock.MatchedBy(func(entry dmntimeline.TimelineEntry) bool {
		now := time.Now().UTC()
		diff := now.Sub(entry.CreatedAt)
		return entry.TweetID == tweet.ID &&
			entry.AuthorID == tweet.UserID &&
			entry.Content == tweet.Content &&
			diff < 5*time.Second
	}), userID).Return(nil)

	err := uc.Exec(context.Background(), tweet, userID)

	assert.NoError(t, err)
	mockTimelineService.AssertExpectations(t)
}

func TestNew_WithNilLogger(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)

	uc := updatetimeline.New(mockTimelineService, nil)

	assert.NotNil(t, uc)

	userID := "user-1"
	tweet := dmntweet.Tweet{
		ID:        "tweet-1",
		UserID:    "author-1",
		Content:   "Hello world!",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	mockTimelineService.On("Update", mock.Anything, mock.Anything, userID).Return(nil)

	err := uc.Exec(context.Background(), tweet, userID)

	assert.NoError(t, err)
	mockTimelineService.AssertExpectations(t)
}
