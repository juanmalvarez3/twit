package daos

import (
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	"time"
)

type TimelineEntryDAO struct {
	UserID    string     `dynamodbav:"user_id" redis:"user_id"`
	SK        string     `dynamodbav:"SK" redis:"SK"`
	TweetID   string     `dynamodbav:"tweet_id" redis:"tweet_id"`
	AuthorID  string     `dynamodbav:"author_id" redis:"author_id"`
	Content   string     `dynamodbav:"content" redis:"content"`
	CreatedAt time.Time  `dynamodbav:"created_at" redis:"created_at"`
	TTL       *time.Time `dynamodbav:"ttl,omitempty" redis:"ttl,omitempty"`
}

func ToTimelineEntryDAO(userID string, entry dmntimeline.TimelineEntry) TimelineEntryDAO {
	content := entry.Content
	if content == "" || userID == "" || entry.TweetID == "" || entry.CreatedAt.IsZero() {
		//TODO: Agregar como validaciones fuera del DAO
		return TimelineEntryDAO{}
	}

	return TimelineEntryDAO{
		UserID:    userID,
		SK:        entry.CreatedAt.Format(time.RFC3339) + "#" + entry.TweetID,
		TweetID:   entry.TweetID,
		AuthorID:  entry.AuthorID,
		Content:   content,
		CreatedAt: entry.CreatedAt,
		TTL:       entry.TTL,
	}
}

func ToDomainTimelineEntry(dao TimelineEntryDAO) dmntimeline.TimelineEntry {
	return dmntimeline.TimelineEntry{
		TweetID:   dao.TweetID,
		AuthorID:  dao.AuthorID,
		Content:   dao.Content,
		CreatedAt: dao.CreatedAt,
		TTL:       dao.TTL,
	}
}
