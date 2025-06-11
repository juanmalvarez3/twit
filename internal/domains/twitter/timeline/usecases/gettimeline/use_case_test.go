package gettimeline_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/gettimeline"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/gettimeline/mocks"
)

func TestExec_Success_WithCacheHit(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockPublisher := new(mocks.Publisher)
	mockFallbackPublisher := new(mocks.FallbackRebuildTimelinePublisherService)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	uc := gettimeline.New(mockTimelineService, mockPublisher, mockFallbackPublisher, mockLogger)

	userID := "user-1"
	now := time.Now().UTC()
	entries := []dmntimeline.TimelineEntry{
		{
			TweetID:   "tweet-1",
			AuthorID:  "author-1",
			Content:   "Hello world!",
			CreatedAt: now,
		},
		{
			TweetID:   "tweet-2",
			AuthorID:  "author-2",
			Content:   "Hello again!",
			CreatedAt: now.Add(-1 * time.Hour),
		},
	}
	timeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: entries,
	}

	mockTimelineService.On("Get", mock.Anything, userID, 30).Return(timeline, true, nil)

	result, err := uc.Exec(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, timeline, result)
	assert.Equal(t, 2, len(result.Entries))
	mockTimelineService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
	mockFallbackPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_Success_WithoutCacheHit(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockPublisher := new(mocks.Publisher)
	mockFallbackPublisher := new(mocks.FallbackRebuildTimelinePublisherService)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	uc := gettimeline.New(mockTimelineService, mockPublisher, mockFallbackPublisher, mockLogger)

	userID := "user-1"
	now := time.Now().UTC()
	entries := []dmntimeline.TimelineEntry{
		{
			TweetID:   "tweet-1",
			AuthorID:  "author-1",
			Content:   "Hello world!",
			CreatedAt: now,
		},
	}
	timeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: entries,
	}

	mockTimelineService.On("Get", mock.Anything, userID, 30).Return(timeline, false, nil)

	mockPublisher.On("Publish", mock.Anything, timeline).Return(nil)

	result, err := uc.Exec(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, timeline, result)
	mockTimelineService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
	mockFallbackPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_EmptyTimeline(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockPublisher := new(mocks.Publisher)
	mockFallbackPublisher := new(mocks.FallbackRebuildTimelinePublisherService)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	uc := gettimeline.New(mockTimelineService, mockPublisher, mockFallbackPublisher, mockLogger)

	userID := "user-1"
	emptyTimeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: []dmntimeline.TimelineEntry{},
	}

	mockTimelineService.On("Get", mock.Anything, userID, 30).Return(emptyTimeline, false, nil)

	mockFallbackPublisher.On("Publish", mock.Anything, userID).Return(nil)

	result, err := uc.Exec(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, dmntimeline.ErrEmptyTimeline, err)
	assert.Equal(t, dmntimeline.Timeline{}, result)
	mockTimelineService.AssertExpectations(t)
	mockFallbackPublisher.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_ServiceError(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockPublisher := new(mocks.Publisher)
	mockFallbackPublisher := new(mocks.FallbackRebuildTimelinePublisherService)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	uc := gettimeline.New(mockTimelineService, mockPublisher, mockFallbackPublisher, mockLogger)

	userID := "user-1"
	expectedErr := errors.New("error obteniendo timeline")
	emptyTimeline := dmntimeline.Timeline{}

	mockTimelineService.On("Get", mock.Anything, userID, 30).Return(emptyTimeline, false, expectedErr)

	result, err := uc.Exec(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, dmntimeline.Timeline{}, result)
	mockTimelineService.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
	mockFallbackPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_FallbackPublisherError(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockPublisher := new(mocks.Publisher)
	mockFallbackPublisher := new(mocks.FallbackRebuildTimelinePublisherService)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	uc := gettimeline.New(mockTimelineService, mockPublisher, mockFallbackPublisher, mockLogger)

	userID := "user-1"
	emptyTimeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: []dmntimeline.TimelineEntry{},
	}
	publisherErr := errors.New("error publicando solicitud de reconstrucción")

	mockTimelineService.On("Get", mock.Anything, userID, 30).Return(emptyTimeline, false, nil)

	mockFallbackPublisher.On("Publish", mock.Anything, userID).Return(publisherErr)

	result, err := uc.Exec(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, dmntimeline.ErrEmptyTimeline, err)
	assert.Equal(t, dmntimeline.Timeline{}, result)
	mockTimelineService.AssertExpectations(t)
	mockFallbackPublisher.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestExec_PublisherError(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockPublisher := new(mocks.Publisher)
	mockFallbackPublisher := new(mocks.FallbackRebuildTimelinePublisherService)
	mockLogger := new(mocks.Logger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	uc := gettimeline.New(mockTimelineService, mockPublisher, mockFallbackPublisher, mockLogger)

	userID := "user-1"
	now := time.Now().UTC()
	entries := []dmntimeline.TimelineEntry{
		{
			TweetID:   "tweet-1",
			AuthorID:  "author-1",
			Content:   "Hello world!",
			CreatedAt: now,
		},
	}
	timeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: entries,
	}
	publisherErr := errors.New("error publicando timeline para lazy caching")

	mockTimelineService.On("Get", mock.Anything, userID, 30).Return(timeline, false, nil)

	mockPublisher.On("Publish", mock.Anything, timeline).Return(publisherErr)

	result, err := uc.Exec(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, timeline, result)
	mockTimelineService.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
	mockFallbackPublisher.AssertNotCalled(t, "Publish")
}

func TestNew_WithNilLogger(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockPublisher := new(mocks.Publisher)
	mockFallbackPublisher := new(mocks.FallbackRebuildTimelinePublisherService)

	assert.Panics(t, func() {
		gettimeline.New(mockTimelineService, mockPublisher, mockFallbackPublisher, nil)
	}, "Se espera un pánico cuando el logger es nil")
}
