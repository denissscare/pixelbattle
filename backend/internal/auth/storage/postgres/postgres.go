package storage

import (
	"database/sql"
	"pixelbattle/internal/auth/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(username, email, passwordHash string) error {
	_, err := r.db.Exec(`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`, username, email, passwordHash)
	return err
}

func (r *Repository) GetUserByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`SELECT id, username, email, password_hash, avatar_url, created_at FROM users WHERE email=$1`, email).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetUserByUsername(username string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`SELECT id, username, email, password_hash, avatar_url, created_at FROM users WHERE username=$1`, username).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
