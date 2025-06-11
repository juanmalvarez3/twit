package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/juanmalvarez3/twit/internal/adapters/queue"

	orchestrateUC "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/orchestratefanout"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic("Error cargando configuración: " + err.Error())
	}

	appLogger, err := logger.New(cfg.Log.Level, cfg.Log.Environment)
	if err != nil {
		panic("Error inicializando logger: " + err.Error())
	}
	defer appLogger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appLogger.Info("Iniciando worker orchestrate-fanout",
		zap.String("env", cfg.Log.Environment),
		zap.String("queue", cfg.SQS.OrchestrateQueue))

	sqsAdapter, err := queue.NewAdapter(cfg)
	if err != nil {
		appLogger.Fatal("Error inicializando adaptador SQS", zap.Error(err))
	}

	orchestrateUseCase := orchestrateUC.Provide(sqsAdapter, cfg, appLogger)

	messageHandler := func(messages []types.Message) error {
		for _, message := range messages {
			appLogger.Info("Procesando mensaje", zap.String("messageId", *message.MessageId))

			var snsMessage struct {
				Type              string                 `json:"Type"`
				MessageId         string                 `json:"MessageId"`
				TopicArn          string                 `json:"TopicArn"`
				Message           string                 `json:"Message"`
				Timestamp         string                 `json:"Timestamp"`
				SignatureVersion  string                 `json:"SignatureVersion"`
				Signature         string                 `json:"Signature"`
				SigningCertURL    string                 `json:"SigningCertURL"`
				UnsubscribeURL    string                 `json:"UnsubscribeURL"`
				MessageAttributes map[string]interface{} `json:"MessageAttributes"`
			}

			appLogger.Debug("Mensaje SQS recibido (raw)", zap.String("body", *message.Body))

			if err := json.Unmarshal([]byte(*message.Body), &snsMessage); err != nil {
				appLogger.Error("Error al deserializar mensaje SNS", zap.Error(err))
				continue
			}

			if snsMessage.Type != "Notification" {
				appLogger.Warn("Tipo de mensaje SNS desconocido", zap.String("type", snsMessage.Type))
				continue
			}

			appLogger.Debug("Mensaje SNS deserializado",
				zap.String("type", snsMessage.Type),
				zap.String("message_id", snsMessage.MessageId),
				zap.String("topic_arn", snsMessage.TopicArn),
				zap.String("raw_message", snsMessage.Message))

			//TODO: Refactor SNS Message
			var rawPayload map[string]interface{}
			if err := json.Unmarshal([]byte(snsMessage.Message), &rawPayload); err != nil {
				appLogger.Error("Error al deserializar mensaje a mapa", zap.Error(err))
				continue
			}

			appLogger.Debug("Contenido deserializado del mensaje", zap.Any("payload", rawPayload))

			var tweetId, userId, content, createdAt string

			if tweetObj, hasTweet := rawPayload["tweet"].(map[string]interface{}); hasTweet {
				if id, ok := tweetObj["id"].(string); ok {
					tweetId = id
				}
				if uid, ok := tweetObj["userId"].(string); ok {
					userId = uid
				}
				if cnt, ok := tweetObj["content"].(string); ok {
					content = cnt
				}
				if cat, ok := tweetObj["createdAt"].(string); ok {
					createdAt = cat
				}

			} else {
				if id, ok := rawPayload["id"].(string); ok {
					tweetId = id
				}
				if uid, ok := rawPayload["userId"].(string); ok {
					userId = uid
				}
				if cnt, ok := tweetObj["content"].(string); ok {
					content = cnt
				}
				if cat, ok := tweetObj["createdAt"].(string); ok {
					createdAt = cat
				}
			}

			appLogger.Info("Datos extraídos del payload",
				zap.String("tweet_id", tweetId),
				zap.String("user_id", userId))
			zap.String("content", content)
			zap.String("created_at", createdAt)

			if tweetId == "" || userId == "" || content == "" || createdAt == "" {
				appLogger.Error("Datos incompletos en el mensaje",
					zap.String("tweet_id", tweetId),
					zap.String("user_id", userId))
				continue
			}

			tweet := domain.Tweet{
				ID:        tweetId,
				UserID:    userId,
				Content:   content,
				CreatedAt: createdAt,
			}

			appLogger.Debug("Tweet recibido para procesar",
				zap.String("tweet_id", tweet.ID),
				zap.String("user_id", tweet.UserID),
				zap.String("content", tweet.Content),
				zap.String("created_at", tweet.CreatedAt))

			appLogger.Info("[ORCHESTRATEFANOUT/main] Iniciando distribución de tweet a timelines de seguidores",
				zap.String("tweet_id", tweet.ID),
				zap.String("user_id", tweet.UserID))

			if err := orchestrateUseCase.Exec(ctx, tweet); err != nil {
				appLogger.Error("Error al distribuir tweet a seguidores",
					zap.Error(err),
					zap.String("tweetId", tweet.ID),
					zap.String("userId", tweet.UserID))
				return err
			}

			appLogger.Info("Tweet distribuido correctamente a seguidores",
				zap.String("tweetId", tweet.ID),
				zap.String("userId", tweet.UserID))
		}
		return nil
	}

	consumer := queue.New(sqsAdapter, cfg.SQS.OrchestrateQueue, messageHandler, appLogger)
	go consumer.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Cerrando worker...")
	cancel()
	appLogger.Info("Worker cerrado correctamente")
}
