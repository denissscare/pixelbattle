package auth

import (
	"errors"
	"pixelbattle/internal/auth/domain"
	storage "pixelbattle/internal/auth/storage/postgres"
	"pixelbattle/pkg/hash"
	jwtutil "pixelbattle/pkg/jwt"
	"pixelbattle/pkg/logger"
)

type Service struct {
	repo       *storage.Repository
	jwtManager *jwtutil.JWTManager
	log        *logger.Logger
}

func NewService(repo *storage.Repository, jwtManager *jwtutil.JWTManager, log *logger.Logger) *Service {
	return &Service{repo: repo, jwtManager: jwtManager, log: log}
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

func (s *Service) Login(emailOrUsername, password string) (*domain.User, string, error) {
	user, err := s.repo.GetUserByEmail(emailOrUsername)
	if err != nil {
		user, err = s.repo.GetUserByUsername(emailOrUsername)
		if err != nil {
			return nil, "", errors.New("user not found")
		}
	}
	if !hash.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}
