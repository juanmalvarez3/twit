package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/juanmalvarez3/twit/pkg/logger"
	"go.uber.org/zap"
	"time"
)

type MessageHandler func([]types.Message) error

type Consumer struct {
	adapter     *Adapter
	queueURL    string
	handler     MessageHandler
	logger      *logger.Logger
	maxMessages int32
	waitTime    int32
}

func New(adapter *Adapter, queueURL string, handler MessageHandler, logger *logger.Logger) *Consumer {
	return &Consumer{
		adapter:     adapter,
		queueURL:    queueURL,
		handler:     handler,
		logger:      logger,
		maxMessages: defaultMaxMessages,
		waitTime:    defaultWaitTime,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	c.logger.Info("Iniciando consumo de mensajes",
		zap.String("queue_url", c.queueURL))

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Deteniendo consumo de mensajes",
				zap.String("queue_url", c.queueURL))
			return
		default:
			messages, err := c.adapter.client.ReceiveMessages(ctx, c.queueURL, c.maxMessages, c.waitTime)
			if err != nil {
				c.logger.Error("Error recibiendo mensajes",
					zap.String("queue_url", c.queueURL),
					zap.Error(err))
				time.Sleep(5 * time.Second) // Ventana de reintentos, debería ser configurable via configs
				continue
			}

			if len(messages) == 0 {
				continue
			}

			if err := c.handler(messages); err != nil {
				c.logger.Error("Error procesando mensajes",
					zap.String("queue_url", c.queueURL),
					zap.Error(err))
				continue
			}

			for _, msg := range messages {
				if err := c.adapter.client.DeleteMessage(ctx, c.queueURL, *msg.ReceiptHandle); err != nil {
					c.logger.Error("Error eliminando mensaje procesado",
						zap.String("queue_url", c.queueURL),
						zap.String("message_id", *msg.MessageId),
						zap.Error(err))
				}
			}
		}
	}
}
