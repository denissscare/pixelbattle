package domain

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	AvatarURL    *string
	CreatedAt    string
}
