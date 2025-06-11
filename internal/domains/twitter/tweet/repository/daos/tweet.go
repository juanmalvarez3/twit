package daos

import (
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"time"
)

type TweetDAO struct {
	ID        string    `json:"id" dynamodbav:"id"`
	UserID    string    `json:"user_id" dynamodbav:"user_id"`
	Content   string    `json:"content" dynamodbav:"content"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"created_at"`
}

func (t *TweetDAO) TableName() string {
	return "tweets"
}

func ToTweetModel(dao TweetDAO) dmntweet.Tweet {
	return dmntweet.Tweet{
		ID:        dao.ID,
		UserID:    dao.UserID,
		Content:   dao.Content,
		CreatedAt: dao.CreatedAt.Format(time.RFC3339),
	}
}

func ToTweetDAOModel(tweetModel dmntweet.Tweet) TweetDAO {
	createdAt, err := time.Parse(time.RFC3339, tweetModel.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}
	return TweetDAO{
		ID:        tweetModel.ID,
		UserID:    tweetModel.UserID,
		Content:   tweetModel.Content,
		CreatedAt: createdAt,
	}
}
