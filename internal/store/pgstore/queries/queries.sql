-- name: GetRoom :one
SELECT
    "id", "subject"
FROM rooms
WHERE id = $1;

-- name: GetRooms :many
SELECT
    "id", "subject"
FROM rooms;

-- name: InsertRoom :one
INSERT INTO rooms
    ( "subject" ) VALUES
    ( $1 )
RETURNING "id";

-- name: GetMessage :one
SELECT
    "id", "room_id", "message", "likes_count", "answered"
FROM messages
WHERE
    id = $1;

-- name: GetRoomMessages :many
SELECT
    "id", "room_id", "message", "likes_count", "answered"
FROM messages
WHERE
    room_id = $1;

-- name: InsertMessage :one
INSERT INTO messages
    ( "room_id", "message" ) VALUES
    ( $1, $2 )
RETURNING "id";

-- name: ReactToMessage :one
UPDATE messages
SET
    likes_count = likes_count + 1
WHERE
    id = $1
RETURNING likes_count;

-- name: RemoveReactionFromMessage :one
UPDATE messages
SET
    likes_count = likes_count - 1
WHERE
    id = $1
RETURNING likes_count;

-- name: MarkMessageAsAnswered :exec
UPDATE messages
SET
    answered = true
WHERE
    id = $1;