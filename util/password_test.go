package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHshPassword(t *testing.T) {
	password := RandomString(6)

	hashpassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashpassword)

	err = ComparePassword(password, hashpassword)
	require.NoError(t, err)

	wrongpass := RandomString(8)
	err = ComparePassword(wrongpass, hashpassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}
