package events

import dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"

const (
	ResourceType = "FOLLOW"

	FollowCreatedEventType EventType = "FOLLOW_CREATED"
)

type EventType string

func (et EventType) String() string {
	return string(et)
}

type Event struct {
	Type     EventType
	Follow   dmntweet.Follow
	Metadata map[string]string
}
