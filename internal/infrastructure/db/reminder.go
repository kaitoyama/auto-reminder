package db

import (
	"context"
	"database/sql"

	"github.com/kaitoyama/kaitoyama-server-template/internal/db"
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
)

type todoReminder struct {
	db *sql.DB
	q  *db.Queries
}

func NewTodoReminder(dbConnection *sql.DB) domain.TodoReminder {
	return &todoReminder{
		db: dbConnection,
		q:  db.New(dbConnection),
	}
}

func (r *todoReminder) convertToDomainTodos(todos []db.Todo) ([]domain.Todo, error) {
	var domainTodos []domain.Todo
	ctx := context.Background()
	for _, todo := range todos {
		users, err := r.q.GetUsersByTodoId(ctx, todo.ID)
		if err != nil {
			return nil, err
		}
		owner, err := r.q.GetUser(ctx, todo.OwnerID)
		domainTodos = append(domainTodos, domain.Todo{
			ID:        int(todo.ID),
			ChannelID: todo.ChannelID, // added mapping for channel_id
			Content:   todo.Content,
			DueAt:     todo.DueAt,
			Users:     users,
			OwnerID:   owner.TraqID,
		})
	}
	return domainTodos, nil
}

func (r *todoReminder) GetTodoInWeek(ctx context.Context) ([]domain.Todo, error) {
	todos, err := r.q.GetTodoInWeek(ctx)
	if err != nil {
		return nil, err
	}

	todoWithUsers, err := r.convertToDomainTodos(todos)
	if err != nil {
		return nil, err
	}

	return todoWithUsers, nil
}

func (r *todoReminder) GetTodoInThreeDays(ctx context.Context) ([]domain.Todo, error) {
	todos, err := r.q.GetTodoInThreeDays(ctx)
	if err != nil {
		return nil, err
	}

	todoWithUsers, err := r.convertToDomainTodos(todos)
	if err != nil {
		return nil, err
	}

	return todoWithUsers, nil
}

func (r *todoReminder) GetTodoInDay(ctx context.Context) ([]domain.Todo, error) {
	todos, err := r.q.GetTodoInDay(ctx)
	if err != nil {
		return nil, err
	}

	todoWithUsers, err := r.convertToDomainTodos(todos)
	if err != nil {
		return nil, err
	}

	return todoWithUsers, nil
}
