package auth

import (
	"errors"
	storage "pixelbattle/internal/auth/storage/postgres"
	"pixelbattle/pkg/hash"
)

type Service struct {
	repo *storage.Repository
}

func NewService(repo *storage.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(username, email, password string) error {
	if _, err := s.repo.GetUserByEmail(email); err == nil {
		return errors.New("user with this email already exists")
	}
	if _, err := s.repo.GetUserByUsername(username); err == nil {
		return errors.New("user with this username already exists")
	}
	hash, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(username, email, hash)
}
