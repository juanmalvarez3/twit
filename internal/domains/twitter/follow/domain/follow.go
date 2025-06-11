package domain

type Follow struct {
	ID         string `json:"id"`
	FollowerID string `json:"followerId"`
	FollowedID string `json:"followedId"`
	CreatedAt  string `json:"createdAt"`
}
