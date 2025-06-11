package domain

import (
	"sort"
	"time"

	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
)

type Timeline struct {
	UserID  string          `json:"user_id"`
	Entries []TimelineEntry `json:"entries"`
}

type TimelineEntry struct {
	TweetID   string     `json:"tweet_id"`
	AuthorID  string     `json:"author_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	TTL       *time.Time `json:"ttl"`
}

func NewTimelineEntryFromTweet(tweet dmntweet.Tweet) TimelineEntry {
	var createdAt time.Time
	if tweet.CreatedAt != "" {
		parsed, err := time.Parse(time.RFC3339, tweet.CreatedAt)
		if err == nil {
			createdAt = parsed
		} else {
			createdAt = time.Now().UTC()
		}
	} else {
		createdAt = time.Now().UTC()
	}

	return TimelineEntry{
		TweetID:   tweet.ID,
		AuthorID:  tweet.UserID,
		Content:   tweet.Content,
		CreatedAt: createdAt,
	}
}

func SortEntriesByTime(entries []TimelineEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt.After(entries[j].CreatedAt)
	})
}
