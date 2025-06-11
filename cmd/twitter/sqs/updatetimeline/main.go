package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/juanmalvarez3/twit/internal/adapters/queue"
	updateTimelineUC "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/updatetimeline"
	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"

	"go.uber.org/zap"
)

type UpdateTimelineRequest struct {
	Tweet  dmntweet.Tweet `json:"tweet"`
	UserID string         `json:"user_id"`
}

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

	appLogger.Info("Iniciando worker update-timeline",
		zap.String("env", cfg.Log.Environment),
		zap.String("queue", cfg.SQS.UpdateTimelineQueue))

	sqsAdapter, err := queue.NewAdapter(cfg)
	if err != nil {
		appLogger.Fatal("Error inicializando adaptador SQS", zap.Error(err))
	}

	updateTimelineUseCase := updateTimelineUC.Provide()

	messageHandler := func(messages []types.Message) error {
		for _, message := range messages {
			messageID := *message.MessageId
			messageBody := *message.Body

			appLogger.Info("Procesando mensaje",
				zap.String("messageId", messageID))
			appLogger.Info("Cuerpo del mensaje SQS",
				zap.String("body", messageBody))

			var updateEvent UpdateTimelineRequest
			if err := json.Unmarshal([]byte(messageBody), &updateEvent); err != nil {
				appLogger.Error("Error al deserializar evento de actualización",
					zap.Error(err),
					zap.String("messageBody", messageBody))
				return err
			}

			appLogger.Info("Tweet deserializado",
				zap.String("tweet_id", updateEvent.Tweet.ID),
				zap.String("user_id", updateEvent.Tweet.UserID),
				zap.String("content", updateEvent.Tweet.Content),
				zap.String("created_at", updateEvent.Tweet.CreatedAt))

			if err := updateTimelineUseCase.Exec(ctx, updateEvent.Tweet, updateEvent.UserID); err != nil {
				appLogger.Error("Error al actualizar timeline",
					zap.Error(err),
					zap.String("userId", updateEvent.UserID),
					zap.String("tweetId", updateEvent.Tweet.ID))
				return err
			}

			appLogger.Info("Timeline actualizado correctamente",
				zap.String("userId", updateEvent.UserID),
				zap.String("tweetId", updateEvent.Tweet.ID))
		}
		return nil
	}

	consumer := queue.New(sqsAdapter, cfg.SQS.UpdateTimelineQueue, messageHandler, appLogger)
	go consumer.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Cerrando worker...")
	cancel()
	appLogger.Info("Worker cerrado correctamente")
}
