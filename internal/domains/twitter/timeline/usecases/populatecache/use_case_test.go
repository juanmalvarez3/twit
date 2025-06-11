package populatecache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/populatecache"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/populatecache/mocks"
)

func TestPopulateCache_Success(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockLogger := new(mocks.Logger)

	uc := populatecache.New(mockTimelineService, mockLogger)

	userID := "user-1"
	inputTimeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: []dmntimeline.TimelineEntry{},
	}

	createdAt1, _ := time.Parse(time.RFC3339, "2025-06-10T22:00:00Z")
	createdAt2, _ := time.Parse(time.RFC3339, "2025-06-10T23:00:00Z")

	dbTimeline := dmntimeline.Timeline{
		UserID: userID,
		Entries: []dmntimeline.TimelineEntry{
			{
				TweetID:   "twt-1",
				AuthorID:  "user-2",
				Content:   "Hello world!",
				CreatedAt: createdAt1,
			},
			{
				TweetID:   "twt-2",
				AuthorID:  "user-3",
				Content:   "Another tweet",
				CreatedAt: createdAt2,
			},
		},
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	mockTimelineService.On("GetFromDB", mock.Anything, userID, 100).Return(dbTimeline, nil)
	mockTimelineService.On("UpdateCache", mock.Anything, dbTimeline).Return(nil)

	err := uc.Exec(context.Background(), inputTimeline)

	assert.NoError(t, err)
	mockTimelineService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestPopulateCache_GetFromDBError(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockLogger := new(mocks.Logger)

	uc := populatecache.New(mockTimelineService, mockLogger)

	userID := "user-1"
	inputTimeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: []dmntimeline.TimelineEntry{},
	}

	expectedErr := errors.New("database error")

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	mockTimelineService.On("GetFromDB", mock.Anything, userID, 100).Return(dmntimeline.Timeline{}, expectedErr)

	err := uc.Exec(context.Background(), inputTimeline)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockTimelineService.AssertExpectations(t)
	mockTimelineService.AssertNotCalled(t, "UpdateCache")
	mockLogger.AssertExpectations(t)
}

func TestPopulateCache_UpdateCacheError(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockLogger := new(mocks.Logger)

	uc := populatecache.New(mockTimelineService, mockLogger)

	userID := "user-1"
	inputTimeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: []dmntimeline.TimelineEntry{},
	}

	createdAt, _ := time.Parse(time.RFC3339, "2025-06-10T22:00:00Z")

	dbTimeline := dmntimeline.Timeline{
		UserID: userID,
		Entries: []dmntimeline.TimelineEntry{
			{
				TweetID:   "twt-1",
				AuthorID:  "user-2",
				Content:   "Hello world!",
				CreatedAt: createdAt,
			},
		},
	}

	expectedErr := errors.New("cache error")

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	mockTimelineService.On("GetFromDB", mock.Anything, userID, 100).Return(dbTimeline, nil)
	mockTimelineService.On("UpdateCache", mock.Anything, dbTimeline).Return(expectedErr)

	err := uc.Exec(context.Background(), inputTimeline)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockTimelineService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestPopulateCache_EmptyTimeline(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockLogger := new(mocks.Logger)

	uc := populatecache.New(mockTimelineService, mockLogger)

	userID := "user-1"
	inputTimeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: []dmntimeline.TimelineEntry{},
	}

	emptyDBTimeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: []dmntimeline.TimelineEntry{},
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()

	mockTimelineService.On("GetFromDB", mock.Anything, userID, 100).Return(emptyDBTimeline, nil)
	mockTimelineService.On("UpdateCache", mock.Anything, emptyDBTimeline).Return(nil)

	err := uc.Exec(context.Background(), inputTimeline)

	assert.NoError(t, err)
	mockTimelineService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestPopulateCache_CanceledContext(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)
	mockLogger := new(mocks.Logger)

	uc := populatecache.New(mockTimelineService, mockLogger)

	userID := "user-1"
	inputTimeline := dmntimeline.Timeline{
		UserID:  userID,
		Entries: []dmntimeline.TimelineEntry{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	contextErr := context.Canceled

	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	mockTimelineService.On("GetFromDB", mock.Anything, userID, 100).Return(dmntimeline.Timeline{}, contextErr)

	err := uc.Exec(ctx, inputTimeline)

	assert.Error(t, err)
	assert.Equal(t, contextErr, err)
	mockTimelineService.AssertExpectations(t)
	mockTimelineService.AssertNotCalled(t, "UpdateCache")
	mockLogger.AssertExpectations(t)
}

func TestNewUseCase_NilLogger(t *testing.T) {
	mockTimelineService := new(mocks.TimelineService)

	assert.Panics(t, func() {
		populatecache.New(mockTimelineService, nil)
	})
}
