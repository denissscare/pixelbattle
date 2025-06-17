package postgres

import (
	"context"
	"database/sql"
	"pixelbattle/internal/pixcelbattle/domain"
)

type PGRepo interface {
	SavePixelHistory(ctx context.Context, p domain.Pixel) error
	GetAllPixelHistory(ctx context.Context) ([]domain.Pixel, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SavePixelHistory(ctx context.Context, p domain.Pixel) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO pixel_history (x, y, color, author, timestamp) VALUES ($1, $2, $3, $4, $5)`,
		p.X, p.Y, p.Color, p.Author, p.Timestamp)
	return err
}

func (r *Repository) GetAllPixelHistory(ctx context.Context) ([]domain.Pixel, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT x, y, color, author, timestamp FROM pixel_history ORDER BY timestamp ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.Pixel
	for rows.Next() {
		var p domain.Pixel
		if err := rows.Scan(&p.X, &p.Y, &p.Color, &p.Author, &p.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, p)
	}
	return history, rows.Err()
}
