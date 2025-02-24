package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simple-banking/util"
	"testing"
)

func createRandomAccount(t *testing.T) Account {

	args := CreateAccountParams{
		Owner:    util.RandomOwnerName(), //randomly generated
		Balance:  int64(util.RandomMoney()),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, args.Owner, acc.Owner)
	return acc
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc.ID, acc2.ID)
}

func TestUpdateAccount(t *testing.T) {
	acc := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      acc.ID,
		Balance: int64(util.RandomMoney()),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, acc.ID, acc2.ID)
	require.Equal(t, args.Balance, acc2.Balance)
}

func TestDeleteAccount(t *testing.T) {
	acc := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc.ID)

	acc2, err2 := testQueries.GetAccount(context.Background(), acc.ID)

	require.NoError(t, err)
	require.Error(t, err2)
	require.EqualError(t, err2, sql.ErrNoRows.Error())
	require.Empty(t, acc2)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := GetAccountsParams{
		Offset: 5,
		Limit:  5,
	}

	accs, err := testQueries.GetAccounts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, accs, 5)

	for _, acc := range accs {
		require.NotEmpty(t, acc.Owner)
	}

}
