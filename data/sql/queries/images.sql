-- name: ImageList :many
select
  *
from
  images
where
  status = 'started'
  OR status = 'succeeded'
  AND owner_id = $1
ORDER BY
  id DESC;

-- name: ImageGet :one
select
  *
from
  images
where
  id = $1
  AND owner_id = $2;

-- name: ImageCreate :one
insert into
  images (provider_id, owner_id, prompt, negative_prompt, status, model)
values
  ($1, $2, $3, $4, $5, $6)
returning
  *;

-- name: ImageUpdate :one
update images
set
  status = $1,
  url = $2
where
  provider_id = $3
returning
  *;
