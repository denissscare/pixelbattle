package auth

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"pixelbattle/internal/auth/domain"
	"pixelbattle/pkg/hash"
	jwtutil "pixelbattle/pkg/jwt"

	"github.com/golang-jwt/jwt/v5"
)

type mockRepo struct {
	getByEmailUser    *domain.User
	getByEmailErr     error
	getByUsernameUser *domain.User
	getByUsernameErr  error
	createUserErr     error
	createID          int
	createUserIDEerr  error
	updateAvatarErr   error
}

func (m *mockRepo) GetUserByEmail(email string) (*domain.User, error) {
	return m.getByEmailUser, m.getByEmailErr
}
func (m *mockRepo) GetUserByUsername(username string) (*domain.User, error) {
	return m.getByUsernameUser, m.getByUsernameErr
}
func (m *mockRepo) CreateUser(username, email, hash string) error {
	return m.createUserErr
}
func (m *mockRepo) CreateUserReturningID(username, email, hash string) (int, error) {
	return m.createID, m.createUserIDEerr
}
func (m *mockRepo) UpdateAvatarURL(id int, url string) error {
	return m.updateAvatarErr
}

type mockJWT struct {
	token    string
	userid   int
	username string
	err      error
}

func (m *mockJWT) GenerateToken(userID int, username string) (string, error) {
	return m.token, m.err
}
func (m *mockJWT) ParseToken(tokenString string) (*jwtutil.Claims, error) {
	return &jwtutil.Claims{
		UserID:           m.userid,
		Username:         m.username,
		RegisteredClaims: jwt.RegisteredClaims{},
	}, m.err
}

type mockS3 struct {
	uploadURL string
	uploadErr error
}

func (m *mockS3) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, objectName string) (string, error) {
	return m.uploadURL, m.uploadErr
}
func (m *mockS3) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	return m.uploadURL, m.uploadErr
}

func newService(repo *mockRepo, jwt *mockJWT, s3 *mockS3) *Service {
	return NewService(repo, jwt, nil, s3)
}

func TestRegisterUserAlreadyExists(t *testing.T) {
	repo := &mockRepo{getByEmailUser: &domain.User{}, getByEmailErr: nil}
	svc := newService(repo, &mockJWT{}, &mockS3{})

	err := svc.Register("test", "mail@mail.com", "pass")
	if err == nil || err.Error() != "user with this email already exists" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterUsernameExists(t *testing.T) {
	repo := &mockRepo{
		getByEmailErr:     errors.New("not found"),
		getByUsernameUser: &domain.User{},
		getByUsernameErr:  nil,
	}
	svc := newService(repo, &mockJWT{}, &mockS3{})

	err := svc.Register("test", "mail@mail.com", "pass")
	if err == nil || err.Error() != "user with this username already exists" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterSuccess(t *testing.T) {
	repo := &mockRepo{
		getByEmailErr:    errors.New("not found"),
		getByUsernameErr: errors.New("not found"),
		createUserErr:    nil,
	}
	svc := newService(repo, &mockJWT{}, &mockS3{})

	err := svc.Register("test", "mail@mail.com", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoginByEmailSuccess(t *testing.T) {
	user := &domain.User{ID: 1, Username: "user", PasswordHash: "hashed"}
	repo := &mockRepo{getByEmailUser: user, getByEmailErr: nil}
	jwt := &mockJWT{token: "token", err: nil}
	svc := newService(repo, jwt, &mockS3{})

	orig := hash.CheckPasswordHash
	hash.CheckPasswordHash = func(password, hash string) bool { return true }
	defer func() { hash.CheckPasswordHash = orig }()

	tokenUser, token, err := svc.Login("mail@mail.com", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenUser.ID != 1 || token != "token" {
		t.Fatalf("wrong result")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	user := &domain.User{ID: 1, Username: "user", PasswordHash: "hashed"}
	repo := &mockRepo{getByEmailUser: user}
	jwt := &mockJWT{token: "token", err: nil}
	svc := newService(repo, jwt, &mockS3{})

	orig := hash.HashPassword
	hash.HashPassword = func(p string) (string, error) { return "hashed", nil }

	defer func() { hash.HashPassword = orig }()

	_, _, err := svc.Login("mail@mail.com", "pass")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterWithIDSuccess(t *testing.T) {
	repo := &mockRepo{
		getByEmailErr:    errors.New("not found"),
		getByUsernameErr: errors.New("not found"),
		createID:         42,
	}
	svc := newService(repo, &mockJWT{}, &mockS3{})

	id, err := svc.RegisterWithID("test", "mail@mail.com", "pass")
	if err != nil || id != 42 {
		t.Fatalf("unexpected: id=%d, err=%v", id, err)
	}
}

func TestUpdateAvatarURLSuccess(t *testing.T) {
	repo := &mockRepo{}
	svc := newService(repo, &mockJWT{}, &mockS3{})
	err := svc.UpdateAvatarURL(1, "url")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUploadAvatarSuccess(t *testing.T) {
	repo := &mockRepo{}
	s3 := &mockS3{}
	svc := newService(repo, &mockJWT{}, s3)

	fileHeader := &multipart.FileHeader{Filename: "avatar.png"}
	err := svc.UploadAvatar(context.Background(), 1, fileHeader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUploadAvatarS3Error(t *testing.T) {
	repo := &mockRepo{}
	s3 := &mockS3{uploadErr: errors.New("s3 failed")}
	svc := newService(repo, &mockJWT{}, s3)

	fileHeader := &multipart.FileHeader{Filename: "avatar.png"}
	err := svc.UploadAvatar(context.Background(), 1, fileHeader)
	if err == nil || err.Error() != "s3 upload error: s3 failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUploadAvatarDBError(t *testing.T) {
	repo := &mockRepo{updateAvatarErr: errors.New("db fail")}
	s3 := &mockS3{}
	svc := newService(repo, &mockJWT{}, s3)

	fileHeader := &multipart.FileHeader{Filename: "avatar.png"}
	err := svc.UploadAvatar(context.Background(), 1, fileHeader)
	if err == nil || err.Error() != "db update error: db fail" {
		t.Fatalf("unexpected error: %v", err)
	}
}
