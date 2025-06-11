package fallbacktimeline

import (
	"context"
	"fmt"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
	dmnoptions "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain/options"
	"go.uber.org/zap"
)

func (u *UseCase) Exec(ctx context.Context, userID string) error {
	logger := zap.L().With(zap.String("userId", userID))
	logger.Info("Iniciando reconstrucción de timeline")

	following, err := u.followsService.GetAllFollowing(ctx, userID)
	if err != nil {
		logger.Error("Error obteniendo seguidos", zap.Error(err))
		return fmt.Errorf("error obteniendo seguidos: %w", err)
	}

	if len(following) == 0 {
		logger.Info("Usuario no sigue a nadie, timeline vacío")
		return nil
	}

	var allTimelineEntries []dmntimeline.TimelineEntry

	for _, follow := range following {
		tweets, _, err := u.tweetService.Search(ctx, dmnoptions.SearchOptions{
			Filters:    dmnoptions.SearchFilters{UserID: &follow},
			Pagination: dmnoptions.SearchPagination{Limit: 10, Offset: 0},
		})
		if err != nil {
			logger.Error("Error obteniendo tweets del usuario seguido",
				zap.Error(err),
				zap.String("followedId", follow))
			return fmt.Errorf("error obteniendo tweets para usuario %s: %w", follow, err)
		}

		for _, tweet := range tweets {
			timelineEntry := dmntimeline.NewTimelineEntryFromTweet(tweet)
			allTimelineEntries = append(allTimelineEntries, timelineEntry)
		}
	}

	if len(allTimelineEntries) == 0 {
		logger.Info("No se encontraron tweets para reconstruir timeline")
		return nil
	}

	dmntimeline.SortEntriesByTime(allTimelineEntries)

	if u.publisher != nil {
		if err := u.publisher.Publish(ctx, userID, allTimelineEntries); err != nil {
			logger.Error("Error publicando timeline reconstruido", zap.Error(err))
			return fmt.Errorf("error publicando timeline reconstruido: %w", err)
		}
		logger.Info("Timeline reconstruido exitosamente", zap.Int("entriesCount", len(allTimelineEntries)))
	}

	return nil
}
