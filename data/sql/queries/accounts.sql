-- name: AccountsList :many
select
  *
from
  accounts;

-- name: AccountGet :one
select
  *
from
  accounts
where
  id = $1;

-- name: AccountGetByUserId :one
select
  *
from
  accounts
where
  user_id = $1;

-- name: AccountCreate :one
insert into
  accounts (user_id, username)
values
  ($1, $2)
returning
  *;

-- name: AccountUpdateUsername :one
update accounts
set
  username = $1
where
  id = $2
returning
  *;
