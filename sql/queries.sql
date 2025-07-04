-- name: CreateUser :one
insert into
  users(
    name,
    email,
    username,
    hashed_password
  )
values
  ($1, $2, $3, $4)
returning
  id;

-- name: FindUserByEmail :one
select
  id,
  hashed_password
from
  users
where
  email = $1
limit
  1;

-- name: FindUserByUsername :one
select
  id,
  hashed_password
from
  users
where
  username = $1
limit
  1;

-- name: FindUserByID :one
select
  id
from
  users
where
  id = $1
limit
  1;
