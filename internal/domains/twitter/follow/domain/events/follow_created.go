package events

type FollowCreatedEvent struct {
	Type     string      `json:"Type"`
	Follow   Follow      `json:"Follow"`
	Metadata interface{} `json:"Metadata"`
}

type Follow struct {
	ID         string `json:"id"`
	FollowerID string `json:"followerId"`
	FollowedID string `json:"followedId"`
	CreatedAt  string `json:"createdAt"`
}
