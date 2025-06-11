package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/juanmalvarez3/twit/internal/adapters/queue"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	rebuildTimelineUC "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/fallbacktimeline"
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

	appLogger.Info("Iniciando worker rebuild-timeline",
		zap.String("env", cfg.Log.Environment),
		zap.String("queue", cfg.SQS.RebuildTimelineQueue))

	sqsAdapter, err := queue.NewAdapter(cfg)
	if err != nil {
		appLogger.Fatal("Error inicializando adaptador SQS", zap.Error(err))
	}

	rebuildTimelineUseCase := rebuildTimelineUC.Provide(
		sqsAdapter,
		cfg,
		appLogger,
	)

	messageHandler := func(messages []types.Message) error {
		for _, message := range messages {
			appLogger.Info("Procesando mensaje", zap.String("messageId", *message.MessageId))

			var populateCacheEvent domain.PopulateCacheEvent
			if err := json.Unmarshal([]byte(*message.Body), &populateCacheEvent); err != nil {
				appLogger.Error("Error al deserializar evento de reconstrucción", zap.Error(err))
				continue
			}

			if err := rebuildTimelineUseCase.Exec(ctx, populateCacheEvent.UserID); err != nil {
				appLogger.Error("Error al reconstruir timeline",
					zap.Error(err),
					zap.String("userId", populateCacheEvent.UserID))
				return err
			}

			appLogger.Info("Timeline reconstruido correctamente",
				zap.String("userId", populateCacheEvent.UserID))
		}
		return nil
	}

	consumer := queue.New(sqsAdapter, cfg.SQS.RebuildTimelineQueue, messageHandler, appLogger)
	go consumer.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Cerrando worker...")
	cancel()
	appLogger.Info("Worker cerrado correctamente")
}
