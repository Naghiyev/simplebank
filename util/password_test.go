package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(6)
	HashPass1, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, HashPass1)

	err = CheckPassword(password, HashPass1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, HashPass1)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPass2, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPass2)
	require.NotEqual(t, HashPass1, hashedPass2)
}
