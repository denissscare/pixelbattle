package hash

import (
	"golang.org/x/crypto/bcrypt"
)

var HashPassword = hashPasswordImpl
var CheckPasswordHash = checkPasswordHashImpl

func hashPasswordImpl(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}
func checkPasswordHashImpl(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
