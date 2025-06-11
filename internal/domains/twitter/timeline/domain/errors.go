package domain

import "errors"

var (
	ErrEmptyTimeline = errors.New("timeline: requested timeline has no entries")
)
