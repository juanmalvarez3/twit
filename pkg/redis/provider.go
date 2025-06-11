package redis

import (
	"github.com/juanmalvarez3/twit/internal/adapters/redis"
	pkgConfig "github.com/juanmalvarez3/twit/pkg/config"
	pkgLogger "github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide() *redis.Client {
	cfg, err := pkgConfig.New()
	if err != nil {
		panic(err)
	}

	pkglogger, err := pkgLogger.ProvideError()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewClient(cfg, pkglogger)
	if err != nil {
		panic(err)
	}

	return redisClient
}
