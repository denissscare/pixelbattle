package hash

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "1488228lalka"
	hash, err := HashPassword(password)
	require.NoError(t, err)
	
	require.Nil(t, bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)))
}
