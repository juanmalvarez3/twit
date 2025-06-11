package services

import (
	"context"
	"fmt"
	"github.com/juanmalvarez3/twit/internal/adapters/sns"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/repository"
	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"
	pkgsns "github.com/juanmalvarez3/twit/pkg/sns"
)

func Provide() Service {
	logs, err := logger.ProvideError()
	if err != nil {
		fmt.Println(err)
	}
	repo := repository.Provide()

	cfg, err := config.New()
	if err != nil {
		fmt.Println("Error cargando configuración:", err)
	}

	awsSnsClient, err := pkgsns.Provide(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	snsClient := sns.NewSNSClient(awsSnsClient, logs)
	publisher := sns.NewTweetSNSPublisher(snsClient, cfg, logs)

	return NewService(repo, publisher, logs)
}
