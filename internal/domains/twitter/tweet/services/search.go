package services

import (
	"context"
	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"go.uber.org/zap"
)

const (
	defaultLimit = 50
	maxLimit     = 100
)

func (s Service) Search(ctx context.Context, opts options.SearchOptions) ([]dmntweet.Tweet, string, error) {
	userID := "<no_user_id>"
	if opts.Filters.UserID != nil {
		userID = *opts.Filters.UserID
	}
	
	s.logger.Debug("Buscando tweets", 
		zap.String("user_id", userID),
		zap.Int("limit_request", opts.Pagination.Limit),
		zap.Int("offset_request", opts.Pagination.Offset),
		zap.String("action", actionSearch),
	)

	if opts.Pagination.Limit == 0 {
		opts.Pagination = opts.Pagination.WithLimit(defaultLimit)
		s.logger.Debug("Ajustando límite a valor por defecto",
			zap.Int("default_limit", defaultLimit),
			zap.String("action", actionSearch),
		)
	}

	if opts.Pagination.Limit > maxLimit {
		opts.Pagination = opts.Pagination.WithLimit(maxLimit)
		s.logger.Debug("Ajustando límite al máximo permitido",
			zap.Int("max_limit", maxLimit),
			zap.String("action", actionSearch),
		)
	}

	if opts.Pagination.Offset <= 0 {
		opts.Pagination = opts.Pagination.WithOffset(0)
	}

	twtList, total, err := s.repository.Search(ctx, userID, opts.Pagination.Limit, "")
	if err != nil {
		s.logger.Error("Error al buscar tweets", 
			zap.String("user_id", userID),
			zap.Int("limit", opts.Pagination.Limit),
			zap.Error(err),
			zap.String("action", actionSearch),
		)
		return nil, "", err
	}

	s.logger.Debug("Tweets encontrados exitosamente", 
		zap.String("user_id", userID),
		zap.Int("count", len(twtList)),
		zap.String("total", total),
		zap.String("action", actionSearch),
	)

	return twtList, total, nil
}
