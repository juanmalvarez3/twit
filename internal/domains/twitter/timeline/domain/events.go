package domain

type PopulateCacheEvent struct {
	UserID string `json:"userId"`
	Source string `json:"source,omitempty"`
}
