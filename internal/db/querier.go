// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateTodo(ctx context.Context, arg CreateTodoParams) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	CreateUserTodoRelation(ctx context.Context, arg CreateUserTodoRelationParams) (sql.Result, error)
	DeleteTodo(ctx context.Context, id int64) (sql.Result, error)
	DeleteUserTodoRelation(ctx context.Context, arg DeleteUserTodoRelationParams) (sql.Result, error)
	DeleteUserTodoRelationByTodoID(ctx context.Context, todoID int64) (sql.Result, error)
	GetTodo(ctx context.Context, id int64) (Todo, error)
	GetTodoInDay(ctx context.Context) ([]Todo, error)
	GetTodoInThreeDays(ctx context.Context) ([]Todo, error)
	GetTodoInWeek(ctx context.Context) ([]Todo, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByTraqId(ctx context.Context, traqID string) (User, error)
	GetUsersByTodoId(ctx context.Context, todoID int64) ([]string, error)
	UpdateTodo(ctx context.Context, arg UpdateTodoParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
