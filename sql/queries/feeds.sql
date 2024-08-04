-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetAllFeeds :many
SELECT * FROM feeds;


-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id = $1;


-- name: GetAllFeedFeedFollowsForAUser :many
SELECT * FROM feed_follows WHERE user_id = $1; 
