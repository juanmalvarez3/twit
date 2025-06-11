package main

import (
	"context"
	"encoding/json"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/events"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/juanmalvarez3/twit/internal/adapters/queue"

	"github.com/juanmalvarez3/twit/internal/adapters/sns"

	orchestrateFanoutUC "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/orchestratefanout"
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

	appLogger.Info("Iniciando worker de procesamiento de tweets creados",
		zap.String("env", cfg.Log.Environment),
		zap.String("queue", cfg.SQS.ProcessFollowQueue))

	sqsAdapter, err := queue.NewAdapter(cfg)
	if err != nil {
		appLogger.Fatal("Error inicializando adaptador SQS", zap.Error(err))
	}

	orchestrateFanoutUseCase := orchestrateFanoutUC.Provide(sqsAdapter, cfg, appLogger)

	messageHandler := func(messages []types.Message) error {
		for _, message := range messages {
			appLogger.Info("Procesando mensaje SNS", zap.String("messageId", *message.MessageId))

			var snsMessage sns.SNSMessage
			if err := json.Unmarshal([]byte(*message.Body), &snsMessage); err != nil {
				appLogger.Error("Error al deserializar mensaje SNS", zap.Error(err))
				continue
			}

			var tweetEvent events.TweetCreatedEvent
			if err := json.Unmarshal([]byte(snsMessage.Message), &tweetEvent); err != nil {
				appLogger.Error("Error al deserializar evento de tweet", zap.Error(err))
				continue
			}

			if err := orchestrateFanoutUseCase.Exec(ctx, tweetEvent.Tweet); err != nil {
				appLogger.Error("Error al procesar tweet creado",
					zap.Error(err),
					zap.String("userId", tweetEvent.Tweet.UserID),
					zap.String("tweetId", tweetEvent.Tweet.ID))
				continue
			}

			appLogger.Info("Tweet procesado correctamente",
				zap.String("userId", tweetEvent.Tweet.UserID),
				zap.String("tweetId", tweetEvent.Tweet.ID))
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
