package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/francisjdev/urlshortener/internal/model"
	"github.com/francisjdev/urlshortener/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresURLRepository struct {
	db *sql.DB
}

func NewPostgresURLRepository(db *sql.DB) *PostgresURLRepository {
	return &PostgresURLRepository{db: db}
}

// compile-time check
var _ repository.URLRepository = (*PostgresURLRepository)(nil)

func (p *PostgresURLRepository) Create(ctx context.Context, url *model.URL) error {
	query := `
        INSERT INTO urls (id, code, long_url, created_at, expires_at, hit_count)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := p.db.ExecContext(ctx, query,
		url.ID,
		url.Code,
		url.LongURL,
		url.CreatedAt,
		url.ExpiresAt,
		url.HitCount,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrCodeAlreadyExists
		}
		return err
	}

	return nil
}
func (p *PostgresURLRepository) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	query := `
		SELECT id, code, long_url, created_at, expires_at, hit_count
		FROM urls
		WHERE code = $1
	`
	url := &model.URL{}
	err := p.db.QueryRowContext(ctx, query, code).Scan(
		&url.ID,
		&url.Code,
		&url.LongURL,
		&url.CreatedAt,
		&url.ExpiresAt,
		&url.HitCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return url, nil
}
func (p *PostgresURLRepository) IncrementHitCount(ctx context.Context, code string) error {
	query := `
		UPDATE urls
		SET hit_count = hit_count + 1
		WHERE code = $1
	`
	res, err := p.db.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
