package domain

type PopulateCacheEvent struct {
	UserID string `json:"user_id"` // Cambiado de "userId" a "user_id" para coincidir con el payload
	Source string `json:"source,omitempty"`
}
