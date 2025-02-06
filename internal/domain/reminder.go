package domain

import (
	"context"
	"errors"
	"time"
)

// HealthStatus represents the health status of a service component
type Todo struct {
	ID        int       `json:"id"`
	ChannelID string    `json:"channel_id"`
	Content   string    `json:"content"`
	DueAt     time.Time `json:"due_at"`
	Users     []string  `json:"users"`
	OwnerID   string    `json:"owner_id"`
}

type TodoCreater interface {
	Create(ctx context.Context, channel_id string, content string, due_at time.Time, owner_id string) (Todo, error)
	AddUser(ctx context.Context, todoID int, userID string) (Todo, error)
	DeleteUser(ctx context.Context, todoID int, userID string) (Todo, error)
	UpdateDueAt(ctx context.Context, todoID int, dueAt time.Time) (Todo, error)
	Delete(ctx context.Context, todoID int) error
}

type TodoReminder interface {
	GetTodoInWeek(ctx context.Context) ([]Todo, error)
	GetTodoInThreeDays(ctx context.Context) ([]Todo, error)
	GetTodoInDay(ctx context.Context) ([]Todo, error)
}

var ErrLastUser = errors.New("cannot remove the last user from todo")
