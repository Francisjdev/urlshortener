package repository

import (
	"context"

	"github.com/francisjdev/urlshortener/internal/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	GetByCode(ctx context.Context, code string) (*model.URL, error)
	IncrementHitCount(ctx context.Context, code string) error
}
