-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, posted_at, feed_id)
VALUES(
$1,
$2,
$3,
$4,
$5,
$6,
$7,
$8
)
RETURNING *;

-- name: GetPosts :many
SELECT * FROM posts
ORDER BY posted_at desc;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
JOIN feeds on feeds.id = posts.feed_id
WHERE feeds.user_id = $1
ORDER BY posts.posted_at
LIMIT $2;
