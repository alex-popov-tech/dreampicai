// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: accounts.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const accountCreate = `-- name: AccountCreate :one
insert into
  accounts (userid, username)
values
  ($1, $2)
returning
  id, userid, username
`

type AccountCreateParams struct {
	Userid   pgtype.UUID
	Username string
}

func (q *Queries) AccountCreate(ctx context.Context, arg AccountCreateParams) (Account, error) {
	row := q.db.QueryRow(ctx, accountCreate, arg.Userid, arg.Username)
	var i Account
	err := row.Scan(&i.ID, &i.Userid, &i.Username)
	return i, err
}

const accountGet = `-- name: AccountGet :one
select
  id, userid, username
from
  accounts
where
  id = $1
`

func (q *Queries) AccountGet(ctx context.Context, id int32) (Account, error) {
	row := q.db.QueryRow(ctx, accountGet, id)
	var i Account
	err := row.Scan(&i.ID, &i.Userid, &i.Username)
	return i, err
}

const accountGetByUserId = `-- name: AccountGetByUserId :one
select
  id, userid, username
from
  accounts
where
  userid = $1
`

func (q *Queries) AccountGetByUserId(ctx context.Context, userid pgtype.UUID) (Account, error) {
	row := q.db.QueryRow(ctx, accountGetByUserId, userid)
	var i Account
	err := row.Scan(&i.ID, &i.Userid, &i.Username)
	return i, err
}

const accountUpdateUsername = `-- name: AccountUpdateUsername :one
update accounts
set
  username = $1
where
  id = $2
returning
  id, userid, username
`

type AccountUpdateUsernameParams struct {
	Username string
	ID       int32
}

func (q *Queries) AccountUpdateUsername(ctx context.Context, arg AccountUpdateUsernameParams) (Account, error) {
	row := q.db.QueryRow(ctx, accountUpdateUsername, arg.Username, arg.ID)
	var i Account
	err := row.Scan(&i.ID, &i.Userid, &i.Username)
	return i, err
}

const accountsList = `-- name: AccountsList :many
select
  id, userid, username
from
  accounts
`

func (q *Queries) AccountsList(ctx context.Context) ([]Account, error) {
	rows, err := q.db.Query(ctx, accountsList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(&i.ID, &i.Userid, &i.Username); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
