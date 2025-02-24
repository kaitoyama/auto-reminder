package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"

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
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return domain.Todo{}, err
	}
	defer tx.Rollback()
	// check owner exists
	user, err := d.q.GetUserByTraqId(ctx, ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Msg("Owner not found")
			// create owner
			userID, err := d.q.CreateUser(ctx, db.CreateUserParams{
				TraqID: ownerID,
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to get user by Traq ID")
				return domain.Todo{}, err
			}
			i, err := userID.LastInsertId()
			if err != nil {
				log.Error().Err(err).Msg("Failed to create user")
				return domain.Todo{}, err
			}
			user = db.User{
				ID:     i,
				TraqID: ownerID,
			}
		} else {
			log.Error().Err(err).Msg("Failed to get last insert ID")
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
		log.Error().Err(err).Msg("Failed to create todo")
		return domain.Todo{}, err
	}

	id64, err := id.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get last insert ID")
		return domain.Todo{}, err
	}

	todo, err := d.q.GetTodo(ctx, id64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get todo")
		return domain.Todo{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
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
		log.Error().Err(err).Msg("Failed to get user by Traq ID")
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
				log.Error().Err(err).Msg("Failed to create user")
				return domain.Todo{}, err
			}
			id, err := userID.LastInsertId()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get last insert ID")
				return domain.Todo{}, err
			}
			user = db.User{
				ID:     id,
				TraqID: traQID,
			}
		} else {
			log.Error().Err(err).Msg("Failed to create user-todo relation")
			return domain.Todo{}, err
		}
	}

	// add user
	_, err = d.q.CreateUserTodoRelation(ctx, db.CreateUserTodoRelationParams{
		UserID: user.ID,
		TodoID: int64(todoID),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return domain.Todo{}, err
	}

	// commit
	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get todo")
		return domain.Todo{}, err
	}

	todo, err := d.q.GetTodo(ctx, int64(todoID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get users by todo ID")
		return domain.Todo{}, err
	}

	users, err := d.q.GetUsersByTodoId(ctx, todo.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return domain.Todo{}, err
	}

	owner, err := d.q.GetUser(ctx, todo.OwnerID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by Traq ID")
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
		log.Error().Err(err).Msg("Failed to get users by todo ID")
		return domain.Todo{}, err
	}

	// 最後の一人の場合はtodoを削除するようにエラーを返す
	users, err := d.q.GetUsersByTodoId(ctx, int64(todoID))
	if err != nil {
		return domain.Todo{}, err
	}
	if len(users) == 1 {
		log.Error().Err(domain.ErrLastUser).Msg("Last user in todo")
		return domain.Todo{}, domain.ErrLastUser
	}

	_, err = d.q.DeleteUserTodoRelation(ctx, db.DeleteUserTodoRelationParams{
		UserID: user.ID,
		TodoID: int64(todoID),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user-todo relation")
		return domain.Todo{}, err
	}

	todo, err := d.q.GetTodo(ctx, int64(todoID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get todo")
		return domain.Todo{}, err
	}

	users, err = d.q.GetUsersByTodoId(ctx, todo.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get users by todo ID")
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
	todo, err := d.q.GetTodo(ctx, int64(todoID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get todo")
		return domain.Todo{}, err
	}

	_, err = d.q.UpdateTodo(ctx, db.UpdateTodoParams{
		ID:        int64(todoID),
		ChannelID: todo.ChannelID,
		Content:   todo.Content,
		DueAt:     dueAt,
		OwnerID:   todo.OwnerID,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update todo")
		return domain.Todo{}, err
	}

	todo, err = d.q.GetTodo(ctx, int64(todoID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get todo")
		return domain.Todo{}, err
	}

	users, err := d.q.GetUsersByTodoId(ctx, todo.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get users by todo ID")
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
	// with transaction
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return err
	}
	defer tx.Rollback()

	_, err = d.q.DeleteUserTodoRelationByTodoID(ctx, int64(todoID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user-todo relation")
		return err
	}

	_, err = d.q.DeleteTodo(ctx, int64(todoID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete todo")
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
	}
	return err
}

func (d *todoCreator) Get(ctx context.Context, todoID int) (domain.Todo, error) {
	todo, err := d.q.GetTodo(ctx, int64(todoID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get todo")
		return domain.Todo{}, err
	}

	users, err := d.q.GetUsersByTodoId(ctx, todo.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get users by todo ID")
		return domain.Todo{}, err
	}

	owner, err := d.q.GetUser(ctx, todo.OwnerID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by Traq ID")
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
