// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
)

const CreateUser = `-- name: CreateUser :execresult
INSERT INTO user (username) VALUES (?)
`

func (q *Queries) CreateUser(ctx context.Context, username string) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateUser, username)
}

const DeleteUser = `-- name: DeleteUser :exec
DELETE FROM user WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, DeleteUser, id)
	return err
}

const GetUser = `-- name: GetUser :one
SELECT id, username FROM user WHERE id = ?
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, GetUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Username)
	return i, err
}

const ListUsers = `-- name: ListUsers :many
SELECT id, username FROM user
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, ListUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Username); err != nil {
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

const UpdateUser = `-- name: UpdateUser :exec
UPDATE user SET username = ? WHERE id = ?
`

type UpdateUserParams struct {
	Username string `json:"username"`
	ID       int32  `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, UpdateUser, arg.Username, arg.ID)
	return err
}
