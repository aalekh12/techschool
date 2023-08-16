-- name: CreateEntry :one
INSERT INTO entries(
    account_id,
    ammount
)VALUES(
    $1,$2
)RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries 
WHERE id = $1 LIMIT 1;

-- name: ListEntry :many
SELECT * FROM entries
WHERE account_id=$1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpateEntry :one
UPDATE entries 
SET ammount=$2
WHERE id=$1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id=$1;