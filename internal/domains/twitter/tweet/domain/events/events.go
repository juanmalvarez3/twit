package events

import dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"

const (
	ResourceType = "TWEET"

	TweetCreatedEventType EventType = "TWEET_CREATED"
)

type EventType string

func (et EventType) String() string {
	return string(et)
}

type Event struct {
	Type     EventType
	Tweet    dmntweet.Tweet
	Metadata map[string]string
}
