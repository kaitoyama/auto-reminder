-- name: CreateTodo :execresult
INSERT INTO todo (channel_id, content, due_at, owner_id)
VALUES (?, ?, ?, ?);
-- name: GetTodo :one
SELECT t.id, t.channel_id, t.content, t.created_at, t.due_at, t.owner_id
FROM todo t
WHERE t.id = ?;

-- name: UpdateTodo :execresult
UPDATE todo
SET channel_id = ?, content = ?, due_at = ?, owner_id = ?
WHERE id = ?;

-- name: DeleteTodo :execresult
DELETE FROM todo
WHERE id = ?;

-- name: GetTodoInWeek :many
SELECT t.id, t.channel_id, t.content, t.created_at, t.due_at, t.owner_id
FROM todo t
WHERE t.due_at >= NOW() + INTERVAL 6 DAY
    AND t.due_at < NOW() + INTERVAL 7 DAY;

-- name: GetTodoInThreeDays :many
SELECT t.id, t.channel_id, t.content, t.created_at, t.due_at, t.owner_id
FROM todo t
WHERE t.due_at >= NOW() + INTERVAL 2 DAY
    AND t.due_at < NOW() + INTERVAL 3 DAY;

-- name: GetTodoInDay :many
SELECT t.id, t.channel_id, t.content, t.created_at, t.due_at, t.owner_id
FROM todo t
WHERE t.due_at >= NOW()
    AND t.due_at < NOW() + INTERVAL 1 DAY;

-- name: CreateUserTodoRelation :execresult
INSERT INTO user_todo_relation (user_id, todo_id)
VALUES (?, ?);

-- name: DeleteUserTodoRelation :execresult
DELETE FROM user_todo_relation
WHERE user_id = ? AND todo_id = ?;

-- name: GetUsersByTodoId :many
SELECT u.traq_id
FROM user u
JOIN user_todo_relation utr ON u.id = utr.user_id
WHERE utr.todo_id = ?;

-- name: CreateUser :execresult
INSERT INTO user (id, traq_id)
VALUES (?, ?);

-- name: GetUser :one
SELECT u.id, u.traq_id
FROM user u
WHERE u.id = ?;

-- name: GetUserByTraqId :one
SELECT u.id, u.traq_id
FROM user u
WHERE u.traq_id = ?;
