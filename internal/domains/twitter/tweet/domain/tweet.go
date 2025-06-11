package domain

import (
	"fmt"
	"strings"
)

type Tweet struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

func (t Tweet) Validate() error {
	if len(t.Content) == 0 {
		return fmt.Errorf("El contenido del tweet no puede estar vacío")
	}
	if len(t.Content) > 280 {
		return fmt.Errorf("El contenido del tweet excede el máximo permitido")
	}
	return nil
}

func (t Tweet) NormalizeContent() string {
	return strings.TrimSpace(t.Content)
}

type TweetCreatedEvent struct {
	Tweet Tweet `json:"tweet"`
}
