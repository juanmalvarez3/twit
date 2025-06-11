package sns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
)

type SNSClient struct {
	client *sns.Client
	logger *logger.Logger
}

func NewSNSClient(client *sns.Client, logger *logger.Logger) *SNSClient {
	return &SNSClient{
		client: client,
		logger: logger,
	}
}

func (c *SNSClient) PublishMessage(ctx context.Context, topicARN string, message interface{}, messageAttributes map[string]string) error {
	c.logger.Debug("Preparando publicación en SNS",
		zap.String("topic_arn", topicARN))

	jsonBytes, err := json.Marshal(message)
	if err != nil {
		c.logger.Error("Error serializando mensaje para SNS",
			zap.String("topic_arn", topicARN),
			zap.Error(err))
		return fmt.Errorf("error serializando mensaje para SNS: %w", err)
	}

	attributes := make(map[string]types.MessageAttributeValue)
	for key, value := range messageAttributes {
		attributes[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	_, err = c.client.Publish(ctx, &sns.PublishInput{
		TopicArn:          aws.String(topicARN),
		Message:           aws.String(string(jsonBytes)),
		MessageAttributes: attributes,
	})

	if err != nil {
		c.logger.Error("Error publicando mensaje en SNS",
			zap.String("topic_arn", topicARN),
			zap.Error(err))
		return fmt.Errorf("error publicando mensaje en SNS: %w", err)
	}

	c.logger.Debug("Mensaje publicado en SNS exitosamente",
		zap.String("topic_arn", topicARN))
	return nil
}

type SNSMessage struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Subject          string `json:"Subject,omitempty"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}
