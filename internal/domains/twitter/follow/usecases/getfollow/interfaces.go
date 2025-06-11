package getfollow

import (
	"context"
	"go.uber.org/zap"
)

type Service interface {
	GetFollowers(ctx context.Context, followedID string) ([]string, error)
	GetFollowing(ctx context.Context, followerID string) ([]string, error)
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}
