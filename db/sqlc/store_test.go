package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTxt(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("before", account1.Balance, account2.Balance)

	// run n concurrent transfer
	n := 3
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("transaction-%d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx,
				TransferTxParams{
					fromAccountID: account1.ID,
					toAccountID:   account2.ID,
					amount:        amount,
				})

			errs <- err
			results <- result
		}()
	}

	//check results

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries

		fromEntry := result.FromEntry
		require.NotZero(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, -amount, fromEntry.Amount)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check account
		fromAccount := result.FromAccount

		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount

		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//check account balance
		fmt.Println("tx", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

		//check final accounts balance

		senderAccount, err := store.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		recieverAccount, err := store.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)
		fmt.Println("after", senderAccount.Balance, recieverAccount.Balance)
		require.Equal(t, account1.Balance-int64(n)*amount, senderAccount.Balance)
		require.Equal(t, account2.Balance+int64(n)*amount, recieverAccount.Balance)

	}

}

func TestTransferTxtWithDeadLock(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("before", account1.Balance, account2.Balance)

	// run n concurrent transfer
	n := 10

	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}
		txName := fmt.Sprintf("transaction-%d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx,
				TransferTxParams{
					fromAccountID: fromAccountId,
					toAccountID:   toAccountId,
					amount:        amount,
				})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
