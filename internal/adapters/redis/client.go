package redis

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

const (
	DefaultTTL = 60 * 60 * time.Second
)

type Client struct {
	client *redis.Client
	logger logger.LoggerInterface
}

func NewClient(cfg *config.Config, logger logger.LoggerInterface) (*Client, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Error("Error conectando a Redis",
			zap.String("error", err.Error()),
			zap.String("host", cfg.Redis.Host),
			zap.String("port", cfg.Redis.Port),
		)
		return nil, err
	}

	logger.Info("Conexión a Redis establecida correctamente",
		zap.String("host", cfg.Redis.Host),
		zap.String("port", cfg.Redis.Port),
	)

	return &Client{
		client: client,
		logger: logger,
	}, nil
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		c.logger.Error("Error obteniendo valor de Redis",
			zap.String("key", key),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	return []byte(val), nil
}

func (c *Client) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	strValue := string(value)
	err := c.client.Set(ctx, key, strValue, ttl).Err()
	if err != nil {
		c.logger.Error("Error estableciendo valor en Redis",
			zap.String("key", key),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	err := c.client.Del(ctx, keys...).Err()
	if err != nil {
		c.logger.Error("Error eliminando clave(s) de Redis",
			zap.Any("keys", keys),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}
