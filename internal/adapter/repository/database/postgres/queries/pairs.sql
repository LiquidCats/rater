-- name: GetAllPairs :many
select * from pairs;

-- name: GetPair :one
select * from pairs where symbol = $1;

-- name: HasPair :one
select exists(select * from pairs where symbol = $1);

-- name: CountPairs :one
select count(*) from pairs;