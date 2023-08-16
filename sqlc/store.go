package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db sql.DB
	Queries
}

func NewStore(db sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: *New(&db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return fmt.Errorf("txerr %v rberr %v", err, rberr)
		}
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Ammount       int64 `json:"ammount"`
}

type TransferTxResult struct {
	Trnasfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Trnasfer, err = q.CreateTransferes(ctx, CreateTransferesParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Ammount:       arg.Ammount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Ammount:   -arg.Ammount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Ammount:   arg.Ammount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Ammount, arg.ToAccountId, arg.Ammount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Ammount, arg.FromAccountId, -arg.Ammount)
		}
		return nil
	})

	return result, err
}

func addMoney(
	context context.Context,
	queries *Queries,
	accountid1 int64,
	ammount1 int64,
	accountId2 int64,
	ammount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = queries.AddAccountBalance(context, AddAccountBalanceParams{
		ID:     accountid1,
		Amount: ammount1,
	})
	if err != nil {
		return
	}

	account2, err = queries.AddAccountBalance(context, AddAccountBalanceParams{
		ID:     accountId2,
		Amount: ammount2,
	})
	return
}
