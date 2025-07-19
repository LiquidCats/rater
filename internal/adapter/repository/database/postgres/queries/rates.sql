-- name: SaveRate :one
insert into rates(price, pair, provider, ts) values ($1, $2, $3, $4) returning *;

-- name: GetRate :one
select * from rates where ts = $1 and pair = $2 limit 1;

-- name: HasRate :one
select exists(select * from rates where ts = $1 and pair = $2);