// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const CreateTodo = `-- name: CreateTodo :execresult
INSERT INTO todo (channel_id, content, due_at, owner_id)
VALUES (?, ?, ?, ?)
`

type CreateTodoParams struct {
	ChannelID string    `json:"channel_id"`
	Content   string    `json:"content"`
	DueAt     time.Time `json:"due_at"`
	OwnerID   int64     `json:"owner_id"`
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateTodo,
		arg.ChannelID,
		arg.Content,
		arg.DueAt,
		arg.OwnerID,
	)
}

const CreateUser = `-- name: CreateUser :execresult
INSERT INTO user (id, traq_id)
VALUES (?, ?)
`

type CreateUserParams struct {
	ID     int64  `json:"id"`
	TraqID string `json:"traq_id"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateUser, arg.ID, arg.TraqID)
}

const CreateUserTodoRelation = `-- name: CreateUserTodoRelation :execresult
INSERT INTO user_todo_relation (user_id, todo_id)
VALUES (?, ?)
`

type CreateUserTodoRelationParams struct {
	UserID int64 `json:"user_id"`
	TodoID int64 `json:"todo_id"`
}

func (q *Queries) CreateUserTodoRelation(ctx context.Context, arg CreateUserTodoRelationParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateUserTodoRelation, arg.UserID, arg.TodoID)
}

const DeleteTodo = `-- name: DeleteTodo :execresult
DELETE FROM todo
WHERE id = ?
`

func (q *Queries) DeleteTodo(ctx context.Context, id int64) (sql.Result, error) {
	return q.db.ExecContext(ctx, DeleteTodo, id)
}

const DeleteUserTodoRelation = `-- name: DeleteUserTodoRelation :execresult
DELETE FROM user_todo_relation
WHERE user_id = ? AND todo_id = ?
`

type DeleteUserTodoRelationParams struct {
	UserID int64 `json:"user_id"`
	TodoID int64 `json:"todo_id"`
}

func (q *Queries) DeleteUserTodoRelation(ctx context.Context, arg DeleteUserTodoRelationParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, DeleteUserTodoRelation, arg.UserID, arg.TodoID)
}

const DeleteUserTodoRelationByTodoID = `-- name: DeleteUserTodoRelationByTodoID :execresult
DELETE FROM user_todo_relation
WHERE todo_id = ?
`

func (q *Queries) DeleteUserTodoRelationByTodoID(ctx context.Context, todoID int64) (sql.Result, error) {
	return q.db.ExecContext(ctx, DeleteUserTodoRelationByTodoID, todoID)
}

const GetTodo = `-- name: GetTodo :one
SELECT t.id, t.channel_id, t.content, t.created_at, t.due_at, t.owner_id
FROM todo t
WHERE t.id = ?
`

func (q *Queries) GetTodo(ctx context.Context, id int64) (Todo, error) {
	row := q.db.QueryRowContext(ctx, GetTodo, id)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.ChannelID,
		&i.Content,
		&i.CreatedAt,
		&i.DueAt,
		&i.OwnerID,
	)
	return i, err
}

const GetTodoInDay = `-- name: GetTodoInDay :many
SELECT t.id, t.channel_id, t.content, t.created_at, t.due_at, t.owner_id
FROM todo t
WHERE t.due_at >= CURDATE()
    AND t.due_at < CURDATE() + INTERVAL 1 DAY
`

func (q *Queries) GetTodoInDay(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, GetTodoInDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Todo{}
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.ChannelID,
			&i.Content,
			&i.CreatedAt,
			&i.DueAt,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetTodoInThreeDays = `-- name: GetTodoInThreeDays :many
SELECT t.id, t.channel_id, t.content, t.created_at, t.due_at, t.owner_id
FROM todo t
WHERE t.due_at >= CURDATE() + INTERVAL 2 DAY
    AND t.due_at < CURDATE() + INTERVAL 3 DAY
`

func (q *Queries) GetTodoInThreeDays(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, GetTodoInThreeDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Todo{}
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.ChannelID,
			&i.Content,
			&i.CreatedAt,
			&i.DueAt,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetTodoInWeek = `-- name: GetTodoInWeek :many
SELECT t.id, t.channel_id, t.content, t.created_at, t.due_at, t.owner_id
FROM todo t
WHERE t.due_at >= CURDATE() + INTERVAL 6 DAY
    AND t.due_at < CURDATE() + INTERVAL 7 DAY
`

func (q *Queries) GetTodoInWeek(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, GetTodoInWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Todo{}
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.ChannelID,
			&i.Content,
			&i.CreatedAt,
			&i.DueAt,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetUser = `-- name: GetUser :one
SELECT u.id, u.traq_id
FROM user u
WHERE u.id = ?
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, GetUser, id)
	var i User
	err := row.Scan(&i.ID, &i.TraqID)
	return i, err
}

const GetUserByTraqId = `-- name: GetUserByTraqId :one
SELECT u.id, u.traq_id
FROM user u
WHERE u.traq_id = ?
`

func (q *Queries) GetUserByTraqId(ctx context.Context, traqID string) (User, error) {
	row := q.db.QueryRowContext(ctx, GetUserByTraqId, traqID)
	var i User
	err := row.Scan(&i.ID, &i.TraqID)
	return i, err
}

const GetUsersByTodoId = `-- name: GetUsersByTodoId :many
SELECT u.traq_id
FROM user u
JOIN user_todo_relation utr ON u.id = utr.user_id
WHERE utr.todo_id = ?
`

func (q *Queries) GetUsersByTodoId(ctx context.Context, todoID int64) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, GetUsersByTodoId, todoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var traq_id string
		if err := rows.Scan(&traq_id); err != nil {
			return nil, err
		}
		items = append(items, traq_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateTodo = `-- name: UpdateTodo :execresult
UPDATE todo
SET channel_id = ?, content = ?, due_at = ?, owner_id = ?
WHERE id = ?
`

type UpdateTodoParams struct {
	ChannelID string    `json:"channel_id"`
	Content   string    `json:"content"`
	DueAt     time.Time `json:"due_at"`
	OwnerID   int64     `json:"owner_id"`
	ID        int64     `json:"id"`
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, UpdateTodo,
		arg.ChannelID,
		arg.Content,
		arg.DueAt,
		arg.OwnerID,
		arg.ID,
	)
}
