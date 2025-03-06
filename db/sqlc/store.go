package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=../mock/store.go -package=db/sqlc simple-bank/db/sqlc Store

// for mocking db functionalites
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// provides all functionalites to slq db
type SqlStore struct {
	*Queries
	db *sql.DB
}

type TransferTxParams struct {
	fromAccountID int64 `json:"from_account_id"`
	toAccountID   int64 `json:"to_account_id"`
	amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}

var txKey = struct{}{}

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *SqlStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.fromAccountID,
			ToAccountID:   arg.toAccountID,
			Amount:        arg.amount,
		})

		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry	1")

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.fromAccountID,
			Amount:    -arg.amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry	2")

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.toAccountID,
			Amount:    arg.amount,
		})

		if err != nil {
			return err
		}

		//account part
		if arg.fromAccountID < arg.toAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, queries, arg.fromAccountID, -arg.amount, arg.toAccountID, arg.amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, queries, arg.toAccountID, arg.amount, arg.fromAccountID, -arg.amount)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return // empty returns are okay as, the return signature has variable names, defined in the function body
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return account1, account2, nil
}
