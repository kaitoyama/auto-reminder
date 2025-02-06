package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/db"
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
)

type todoCreator struct {
	db *sql.DB
	q  *db.Queries
}

func NewTodoCreator(dbConnection *sql.DB) domain.TodoCreater {
	return &todoCreator{
		db: dbConnection,
		q:  db.New(dbConnection),
	}
}

func (d *todoCreator) Create(ctx context.Context, channelID string, content string, dueAt time.Time, ownerID string) (domain.Todo, error) {
	// check owner exists
	user, err := d.q.GetUserByTraqId(ctx, ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			// create owner
			userID, err := d.q.CreateUser(ctx, db.CreateUserParams{
				TraqID: ownerID,
			})
			if err != nil {
				return domain.Todo{}, err
			}
			i, err := userID.LastInsertId()
			if err != nil {
				return domain.Todo{}, err
			}
			user = db.User{
				ID:     i,
				TraqID: ownerID,
			}
		} else {
			return domain.Todo{}, err
		}
	}

	id, err := d.q.CreateTodo(ctx, db.CreateTodoParams{
		ChannelID: channelID,
		Content:   content,
		DueAt:     dueAt,
		OwnerID:   user.ID,
	})
	if err != nil {
		return domain.Todo{}, err
	}

	id64, err := id.LastInsertId()
	if err != nil {
		return domain.Todo{}, err
	}

	todo, err := d.q.GetTodo(ctx, id64)
	if err != nil {
		return domain.Todo{}, err
	}

	return domain.Todo{
		ID:        int(todo.ID),
		ChannelID: todo.ChannelID,
		Content:   todo.Content,
		DueAt:     todo.DueAt,
		OwnerID:   ownerID,
	}, nil

}

func (d *todoCreator) AddUser(ctx context.Context, todoID int, traQID string) (domain.Todo, error) {
	// with transaction
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Todo{}, err
	}
	defer tx.Rollback()

	// search user
	user, err := d.q.GetUserByTraqId(ctx, traQID)
	if err != nil {
		if err == sql.ErrNoRows {
			userID, err := d.q.CreateUser(ctx, db.CreateUserParams{
				TraqID: traQID,
			})
			if err != nil {
				return domain.Todo{}, err
			}
			id, err := userID.LastInsertId()
			if err != nil {
				return domain.Todo{}, err
			}
			user = db.User{
				ID:     id,
				TraqID: traQID,
			}
		} else {
			return domain.Todo{}, err
		}
	}

	// add user
	_, err = d.q.CreateUserTodoRelation(ctx, db.CreateUserTodoRelationParams{
		UserID: user.ID,
		TodoID: int64(todoID),
	})
	if err != nil {
		return domain.Todo{}, err
	}

	// commit
	err = tx.Commit()
	if err != nil {
		return domain.Todo{}, err
	}

	todo, err := d.q.GetTodo(ctx, int64(todoID))
	if err != nil {
		return domain.Todo{}, err
	}

	users, err := d.q.GetUsersByTodoId(ctx, todo.ID)
	if err != nil {
		return domain.Todo{}, err
	}

	owner, err := d.q.GetUser(ctx, todo.OwnerID)
	if err != nil {
		return domain.Todo{}, err
	}

	return domain.Todo{
		ID:        int(todo.ID),
		ChannelID: todo.ChannelID,
		Content:   todo.Content,
		DueAt:     todo.DueAt,
		Users:     users,
		OwnerID:   owner.TraqID,
	}, nil
}

func (d *todoCreator) DeleteUser(ctx context.Context, todoID int, userID string) (domain.Todo, error) {
	user, err := d.q.GetUserByTraqId(ctx, userID)
	if err != nil {
		return domain.Todo{}, err
	}

	// 最後の一人の場合はtodoを削除するようにエラーを返す
	users, err := d.q.GetUsersByTodoId(ctx, int64(todoID))
	if err != nil {
		return domain.Todo{}, err
	}
	if len(users) == 1 {
		return domain.Todo{}, domain.ErrLastUser
	}

	_, err = d.q.DeleteUserTodoRelation(ctx, db.DeleteUserTodoRelationParams{
		UserID: user.ID,
		TodoID: int64(todoID),
	})
	if err != nil {
		return domain.Todo{}, err
	}

	todo, err := d.q.GetTodo(ctx, int64(todoID))
	if err != nil {
		return domain.Todo{}, err
	}

	users, err = d.q.GetUsersByTodoId(ctx, todo.ID)
	if err != nil {
		return domain.Todo{}, err
	}

	return domain.Todo{
		ID:        int(todo.ID),
		ChannelID: todo.ChannelID,
		Content:   todo.Content,
		DueAt:     todo.DueAt,
		Users:     users,
	}, nil
}

func (d *todoCreator) UpdateDueAt(ctx context.Context, todoID int, dueAt time.Time) (domain.Todo, error) {
	_, err := d.q.UpdateTodo(ctx, db.UpdateTodoParams{
		ID:    int64(todoID),
		DueAt: dueAt,
	})
	if err != nil {
		return domain.Todo{}, err
	}

	todo, err := d.q.GetTodo(ctx, int64(todoID))
	if err != nil {
		return domain.Todo{}, err
	}

	users, err := d.q.GetUsersByTodoId(ctx, todo.ID)
	if err != nil {
		return domain.Todo{}, err
	}

	return domain.Todo{
		ID:        int(todo.ID),
		ChannelID: todo.ChannelID,
		Content:   todo.Content,
		DueAt:     todo.DueAt,
		Users:     users,
	}, nil
}

func (d *todoCreator) Delete(ctx context.Context, todoID int) error {
	_, err := d.q.DeleteTodo(ctx, int64(todoID))
	return err
}
