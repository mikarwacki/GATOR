-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, 
	created_at,
	updated_at,
	user_id,
	feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, users.name as username, feeds.name as feed_name FROM feed_follows 
JOIN users on users.id = feed_follows.user_id
JOIN feeds on feeds.id = feed_follows.feed_id
WHERE users.name = $1;