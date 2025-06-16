package storage

import (
	"database/sql"
	"pixelbattle/internal/auth/domain"
)

type UserRepo interface {
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	CreateUser(username, email, hash string) error
	CreateUserReturningID(username, email, hash string) (int, error)
	UpdateAvatarURL(id int, url string) error
}

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

func (r *Repository) UpdateAvatarURL(userID int, avatarURL string) error {
	_, err := r.db.Exec(`UPDATE users SET avatar_url = $1 WHERE id = $2`, avatarURL, userID)
	return err
}

func (r *Repository) CreateUserReturningID(username, email, password string) (int, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		username, email, password,
	).Scan(&id)
	return id, err
}
