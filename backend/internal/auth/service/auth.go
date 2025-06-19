package auth

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"pixelbattle/internal/auth/domain"
	storage "pixelbattle/internal/auth/storage/postgres"
	"pixelbattle/internal/s3"
	"pixelbattle/pkg/hash"
	jwtutil "pixelbattle/pkg/jwt"
	"pixelbattle/pkg/logger"
	"strings"
)

type Service struct {
	repo       storage.UserRepo
	jwtManager jwtutil.JWT
	log        *logger.Logger
	s3Client   s3.Uploader
}

func NewService(repo storage.UserRepo, jwtManager jwtutil.JWT, log *logger.Logger, s3Client s3.Uploader) *Service {
	return &Service{repo: repo, jwtManager: jwtManager, log: log, s3Client: s3Client}
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

func (s *Service) RegisterWithID(username, email, password string) (int, error) {
	if _, err := s.repo.GetUserByEmail(email); err == nil {
		return 0, errors.New("почта занята")
	}

	if _, err := s.repo.GetUserByUsername(username); err == nil {
		return 0, errors.New("логин занят")
	}

	hash, err := hash.HashPassword(password)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateUserReturningID(username, email, hash)
}

func (s *Service) UpdateAvatarURL(userID int, url string) error {
	return s.repo.UpdateAvatarURL(userID, url)
}

func (s *Service) UploadAvatar(ctx context.Context, userID int, fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	objectName := fmt.Sprintf("%d%s", userID, ext)

	_, err := s.s3Client.UploadFile(ctx, fileHeader, objectName)
	if err != nil {
		return fmt.Errorf("s3 upload error: %w", err)
	}

	if err := s.repo.UpdateAvatarURL(userID, objectName); err != nil {
		return fmt.Errorf("db update error: %w", err)
	}
	return nil
}

func (s *Service) UpdateEmail(userID int, email string) error {
	if _, err := s.repo.GetUserByEmail(email); err == nil {
		return errors.New("почта уже используется")
	}
	return s.repo.UpdateEmail(userID, email)
}
