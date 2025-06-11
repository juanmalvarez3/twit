package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
)

type SQSClient struct {
	client *sqs.Client
	logger *logger.Logger
}

func NewSQSClient(client *sqs.Client, logger *logger.Logger) *SQSClient {
	return &SQSClient{
		client: client,
		logger: logger,
	}
}

func (c *SQSClient) Publish(ctx context.Context, queueURL string, payload any) error {
	c.logger.Debug("Preparando envío de mensaje a SQS",
		zap.String("queue_url", queueURL))

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		c.logger.Error("Error serializando payload para SQS",
			zap.Error(err))
		return fmt.Errorf("error serializando payload: %w", err)
	}

	_, err = c.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(jsonBytes)),
	})

	if err != nil {
		c.logger.Error("Error enviando mensaje a SQS",
			zap.String("queue_url", queueURL),
			zap.Error(err))
		return fmt.Errorf("error enviando mensaje a SQS: %w", err)
	}

	c.logger.Debug("Mensaje enviado a SQS exitosamente",
		zap.String("queue_url", queueURL))
	return nil
}

func (c *SQSClient) ReceiveMessages(ctx context.Context, queueURL string, maxMessages int32, waitTimeSeconds int32) ([]types.Message, error) {
	c.logger.Debug("Recibiendo mensajes de SQS",
		zap.String("queue_url", queueURL),
		zap.Int32("max_messages", maxMessages),
		zap.Int32("wait_time_seconds", waitTimeSeconds))

	result, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTimeSeconds,
	})

	if err != nil {
		c.logger.Error("Error recibiendo mensajes de SQS",
			zap.String("queue_url", queueURL),
			zap.Error(err))
		return nil, fmt.Errorf("error recibiendo mensajes de SQS: %w", err)
	}

	c.logger.Debug("Mensajes recibidos de SQS",
		zap.String("queue_url", queueURL),
		zap.Int("message_count", len(result.Messages)))
	return result.Messages, nil
}

func (c *SQSClient) DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error {
	c.logger.Debug("Eliminando mensaje de SQS",
		zap.String("queue_url", queueURL),
		zap.String("receipt_handle", receiptHandle))

	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})

	if err != nil {
		c.logger.Error("Error eliminando mensaje de SQS",
			zap.String("queue_url", queueURL),
			zap.String("receipt_handle", receiptHandle),
			zap.Error(err))
		return fmt.Errorf("error eliminando mensaje de SQS: %w", err)
	}

	c.logger.Debug("Mensaje eliminado de SQS exitosamente",
		zap.String("queue_url", queueURL),
		zap.String("receipt_handle", receiptHandle))
	return nil
}

func (c *SQSClient) Send(ctx context.Context, queueURL string, payload any) error {
	return c.Publish(ctx, queueURL, payload)
}
