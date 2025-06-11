package daos

import (
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"time"
)

type FollowDAO struct {
	ID         string    `json:"id" dynamodbav:"id"`
	FollowerID string    `json:"follower_id" dynamodbav:"follower_id"`
	FollowedID string    `json:"followed_id" dynamodbav:"followed_id"`
	CreatedAt  time.Time `json:"created_at" dynamodbav:"created_at"`
}

func (f *FollowDAO) TableName() string {
	return "follows"
}

func ToFollowModel(dao FollowDAO) dmnfollow.Follow {
	follow := dmnfollow.Follow{
		ID:         dao.ID,
		FollowerID: dao.FollowerID,
		FollowedID: dao.FollowedID,
		CreatedAt:  dao.CreatedAt.Format(time.RFC3339),
	}

	return follow
}

func ToFollowDAOModel(followModel dmnfollow.Follow) FollowDAO {
	createdAt := time.Now().UTC()
	if followModel.CreatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, followModel.CreatedAt)
		if err == nil {
			createdAt = parsedTime
		}
	}

	follow := FollowDAO{
		ID:         followModel.ID,
		FollowerID: followModel.FollowerID,
		FollowedID: followModel.FollowedID,
		CreatedAt:  createdAt,
	}

	return follow
}
