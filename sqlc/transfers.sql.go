// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: transfers.sql

package db

import (
	"context"
)

const createTransferes = `-- name: CreateTransferes :one
INSERT INTO transfers(
    from_account_id,
    to_account_id,
    ammount
)VALUES(
    $1,$2,$3
)RETURNING id, from_account_id, to_account_id, ammount, created_at
`

type CreateTransferesParams struct {
	FromAccountID int64`json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Ammount       int64         `json:"ammount"`
}

func (q *Queries) CreateTransferes(ctx context.Context, arg CreateTransferesParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransferes, arg.FromAccountID, arg.ToAccountID, arg.Ammount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Ammount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, ammount, created_at FROM transfers
WHERE id=$1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Ammount,
		&i.CreatedAt,
	)
	return i, err
}

const listTransferes = `-- name: ListTransferes :many
SELECT id, from_account_id, to_account_id, ammount, created_at FROM transfers
WHERE
      from_account_id =$1 OR
      to_account_id=$2
      ORDER BY id
LIMIT $3
OFFSET $4
`

type ListTransferesParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Limit         int32         `json:"limit"`
	Offset        int32         `json:"offset"`
}

func (q *Queries) ListTransferes(ctx context.Context, arg ListTransferesParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransferes,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Ammount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
