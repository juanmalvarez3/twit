package service

import (
	"context"
	"encoding/json"
	dmntimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/domain"
)

func (s Service) UpdateCache(ctx context.Context, timeline dmntimeline.Timeline) error {
	cacheKey := "timeline:" + timeline.UserID
	timelineToSave, err := json.Marshal(timeline)
	if err != nil {
		return err
	}

	return s.timelineRepo.SetCache(ctx, cacheKey, timelineToSave)
}
