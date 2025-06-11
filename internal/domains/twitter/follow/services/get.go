package services

import (
	"context"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"go.uber.org/zap"
)

func (s Service) GetFollowers(ctx context.Context, followedID string) ([]string, error) {
	s.logger.Debug("Servicio: Obteniendo seguidores", zap.String("followed_id", followedID))

	return s.repository.GetFollowers(ctx, followedID)
}

func (s Service) GetAllFollowing(ctx context.Context, followerID string) ([]string, error) {
	s.logger.Debug("Servicio: Obteniendo usuarios seguidos", zap.String("follower_id", followerID))

	return s.repository.GetFollowing(ctx, followerID)
}

func (s Service) GetFollowing(ctx context.Context, followerID string) ([]string, error) {
	s.logger.Debug("Servicio: Obteniendo usuarios seguidos", zap.String("follower_id", followerID))

	return s.repository.GetFollowing(ctx, followerID)
}

func (s Service) Get(ctx context.Context, followID string) (dmnfollow.Follow, error) {
	s.logger.Debug("Servicio: Obteniendo relacion de usuarios", zap.String("follow_id", followID))

	return s.repository.Get(ctx, followID)
}
