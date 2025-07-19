-- name: GetAllProviders :many
select * from providers;

-- name: GetProvider :one
select * from providers where name = $1;