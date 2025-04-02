package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple-banking/util"
	"testing"
)

func createRandomUser(t *testing.T) User {

	args := CreateUserParams{
		Username:       util.RandomOwnerName(), //randomly generated
		HashedPassword: "int64(util.RandomMoney())",
		FullName:       util.RandomOwnerName(),
		Email:          util.RandomEmailAddress(),
	}

	usr, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, usr)
	require.Equal(t, args.FullName, usr.FullName)
	require.NotZero(t, usr.CreatedAt)
	require.True(t, usr.PasswordChangedAt.IsZero())
	return usr
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	usr := createRandomUser(t)
	usr2, err := testQueries.GetUser(context.Background(), usr.Username)
	require.NoError(t, err)
	require.NotEmpty(t, usr2)
	require.Equal(t, usr.Username, usr2.Username)
}
