package events

import dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"

type TweetCreatedEvent struct {
	Tweet dmntweet.Tweet `json:"tweet"`
}

func NewTweetCreatedEvent(tweet dmntweet.Tweet) TweetCreatedEvent {
	return TweetCreatedEvent{
		Tweet: tweet,
	}
}
