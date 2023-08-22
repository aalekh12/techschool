package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/samplebank/util"
)

func CreateUserAccount(t *testing.T) User {
	hashpassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.GenerateUser(),
		HashedPassword: hashpassword,
		FullName:       util.GenerateUser(),
		Email:          util.RandomEmail(),
	}

	user, err := testQuries.CreateUser(context.Background(), arg)
	if err != nil {
		fmt.Println("Error in Creating Account", err)
	}

	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChnagedAt.IsZero())

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUserAccount(t *testing.T) {
	CreateUserAccount(t)
}

func TestGetuses(t *testing.T) {
	user1 := CreateUserAccount(t)
	user2, err := testQuries.GetUsers(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.CreatedAt, user1.CreatedAt, time.Second)

}
