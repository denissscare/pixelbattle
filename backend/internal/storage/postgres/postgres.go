package postgres

import (
	"database/sql"
	"fmt"
	"pixelbattle/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewStorage(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName, cfg.Postgres.SSLMode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
