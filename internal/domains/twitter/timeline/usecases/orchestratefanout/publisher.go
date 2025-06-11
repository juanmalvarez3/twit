package orchestratefanout

import (
	"context"

	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
)

type SQSClient interface {
	Send(ctx context.Context, queueURL string, payload any) error
}

type UpdateTimelinePayload struct {
	Tweet struct {
		ID        string `json:"id"`
		UserID    string `json:"userId"`
		Content   string `json:"content"`
		CreatedAt string `json:"createdAt"`
	} `json:"tweet"`
	TimelineID string `json:"user_id"`
}

type SQSPublisher struct {
	Client   SQSClient
	QueueURL string
}

func UpdateTimelinePublisher(client SQSClient, queueURL string) *SQSPublisher {
	return &SQSPublisher{
		Client:   client,
		QueueURL: queueURL,
	}
}

func (p *SQSPublisher) Publish(ctx context.Context, tweet dmntweet.Tweet, timelineID string) error {
	var payload UpdateTimelinePayload

	payload.Tweet.ID = tweet.ID
	payload.Tweet.UserID = tweet.UserID
	payload.Tweet.Content = tweet.Content
	payload.Tweet.CreatedAt = tweet.CreatedAt
	payload.TimelineID = timelineID

	return p.Client.Send(ctx, p.QueueURL, payload)
}
