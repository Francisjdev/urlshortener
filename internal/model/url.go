package model

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID        uuid.UUID
	Code      string
	LongURL   string
	CreatedAt time.Time
	ExpiresAt *time.Time
	HitCount  int
}
