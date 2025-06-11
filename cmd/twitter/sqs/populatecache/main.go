package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/juanmalvarez3/twit/internal/adapters/queue"

	ucpopulatecache "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/populatecache"
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

	appLogger.Info("Iniciando worker populate-cache",
		zap.String("env", cfg.Log.Environment),
		zap.String("queue", cfg.SQS.OrchestrateQueue))

	sqsAdapter, err := queue.NewAdapter(cfg)
	if err != nil {
		appLogger.Fatal("Error inicializando adaptador SQS", zap.Error(err))
	}

	populateCacheUC := ucpopulatecache.Provide(appLogger)

	messageHandler := func(messages []types.Message) error {
		for _, message := range messages {
			messageID := *message.MessageId
			messageBody := *message.Body

			appLogger.Info("Procesando mensaje",
				zap.String("messageId", messageID))
			appLogger.Info("Cuerpo del mensaje SQS",
				zap.String("body", messageBody))

			var timeline dmntimeline.Timeline
			if err := json.Unmarshal([]byte(messageBody), &timeline); err != nil {
				appLogger.Error("Error al deserializar timeline",
					zap.Error(err),
					zap.String("messageBody", messageBody))
				return err
			}

			appLogger.Info("Timeline deserializado",
				zap.String("user_id", timeline.UserID))

			if err := populateCacheUC.Exec(ctx, timeline); err != nil {
				appLogger.Error("Error al actualizar timeline",
					zap.Error(err))
				return err
			}

			appLogger.Info("Timeline actualizado correctamente",
				zap.String("userId", timeline.UserID),
				zap.Int("entries_count", len(timeline.Entries)))
		}
		return nil
	}

	consumer := queue.New(sqsAdapter, cfg.SQS.PopulateCacheQueue, messageHandler, appLogger)
	go consumer.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Cerrando worker...")
	cancel()
	appLogger.Info("Worker cerrado correctamente")
}
