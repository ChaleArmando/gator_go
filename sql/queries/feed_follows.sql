-- name: CreateFeedFollow :one
WITH feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1, $2, $3, $4, $5
    )
    RETURNING *
)
SELECT ff.*, ur.name AS user_name, fd.name AS feed_name
FROM feed_follow ff
INNER JOIN users ur ON ff.user_id = ur.id
INNER JOIN feeds fd ON ff.feed_id = fd.id;

-- name: GetFeedFollowsForUser :many
SELECT ff.*, ur.name AS user_name, fd.name AS feed_name
FROM feed_follows ff
INNER JOIN users ur ON ff.user_id = ur.id
INNER JOIN feeds fd ON ff.feed_id = fd.id
WHERE ur.name = $1;