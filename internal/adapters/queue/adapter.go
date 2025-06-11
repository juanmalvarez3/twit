package queue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	appConfig "github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

type Adapter struct {
	client           *SQSClient
	orchestrateQueue string
	updateQueue      string
	processQueue     string
	populateQueue    string
	rebuildQueue     string
}

func NewAdapter(cfg *appConfig.Config) (*Adapter, error) {
	appLogger, err := logger.New(cfg.Log.Level, cfg.Log.Environment)
	if err != nil {
		return nil, fmt.Errorf("error inicializando logger para SQS: %w", err)
	}

	credProvider := credentials.NewStaticCredentialsProvider(
		cfg.AWS.AccessKey,
		cfg.AWS.SecretKey,
		"",
	)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           cfg.AWS.Endpoint,
			SigningRegion: cfg.AWS.Region,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credProvider),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("error cargando configuración AWS para SQS: %w", err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)
	client := NewSQSClient(sqsClient, appLogger)

	return &Adapter{
		client:           client,
		orchestrateQueue: cfg.SQS.OrchestrateQueue,
		updateQueue:      cfg.SQS.UpdateTimelineQueue,
		processQueue:     cfg.SQS.ProcessFollowQueue,
		populateQueue:    cfg.SQS.PopulateCacheQueue,
		rebuildQueue:     cfg.SQS.RebuildTimelineQueue,
	}, nil
}

func (a *Adapter) Send(ctx context.Context, queueURL string, payload any) error {
	return a.client.Send(ctx, queueURL, payload)
}
