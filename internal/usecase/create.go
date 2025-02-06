package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/rs/zerolog/log"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

type CreateUsecase struct {
	traQWSBot *traqwsbot.Bot
	creator   domain.TodoCreater
}

func NewCreateUsecase(creator domain.TodoCreater, bot *traqwsbot.Bot) *CreateUsecase {
	return &CreateUsecase{
		creator:   creator,
		traQWSBot: bot,
	}
}

func (u *CreateUsecase) Create(ctx context.Context, channelID string, content string, dueAt time.Time, userTraQIDs []string, ownerID string) error {
	todo, err := u.creator.Create(ctx, channelID, content, dueAt, ownerID)
	if err != nil {
		return err
	}

	for _, userTraQID := range userTraQIDs {
		_, err := u.creator.AddUser(ctx, todo.ID, userTraQID)
		if err != nil {
			return err
		}
	}

	_, _, err = u.traQWSBot.API().MessageApi.PostMessage(ctx, channelID).PostMessageRequest(
		traq.PostMessageRequest{
			Content: fmt.Sprintf(` 新しいリマインドを作成しました!  id: %d`, todo.ID),
		},
	).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (u *CreateUsecase) AddUser(ctx context.Context, todoID int, userTraQIDs []string) error {
	var failed []string
	var todo *domain.Todo
	for _, userID := range userTraQIDs {
		t, err := u.creator.AddUser(ctx, todoID, userID)
		if err != nil {
			failed = append(failed, userID)
		} else {
			// use the first success for the todo info (ChannelID, etc.)
			if todo == nil {
				todo = &t
			}
		}
	}

	// If none of the users could be added, return an error.
	if todo == nil {
		return fmt.Errorf("ユーザーの追加に失敗しました: %v", failed)
	}

	message := fmt.Sprintf("リマインドにユーザーを追加しました! id: %d", todo.ID)
	if len(failed) > 0 {
		message += fmt.Sprintf(" なお、%vは追加できませんでした", failed)
	}

	_, _, err := u.traQWSBot.API().MessageApi.PostMessage(ctx, todo.ChannelID).PostMessageRequest(
		traq.PostMessageRequest{
			Content: message,
		},
	).Execute()
	if err != nil {
		return err
	}
	return nil
}

func (u *CreateUsecase) DeleteUser(ctx context.Context, todoID int, userTraQID string) error {
	todo, err := u.creator.DeleteUser(ctx, todoID, userTraQID)
	if err != nil {
		if err == domain.ErrLastUser {
			_, _, err = u.traQWSBot.API().MessageApi.PostMessage(ctx, todo.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: fmt.Sprintf(`最後のユーザーを削除することはできません!  id: %d`, todo.ID),
				},
			).Execute()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	_, _, err = u.traQWSBot.API().MessageApi.PostMessage(ctx, todo.ChannelID).PostMessageRequest(
		traq.PostMessageRequest{
			Content: fmt.Sprintf(` リマインドからユーザーを削除しました!  id: %d`, todo.ID),
		},
	).Execute()
	if err != nil {
		return err
	}

	return nil

}

func (u *CreateUsecase) UpdateDueAt(ctx context.Context, todoID int, dueAt time.Time) error {
	todo, err := u.creator.UpdateDueAt(ctx, todoID, dueAt)
	if err != nil {
		return err
	}

	_, _, err = u.traQWSBot.API().MessageApi.PostMessage(ctx, todo.ChannelID).PostMessageRequest(
		traq.PostMessageRequest{
			Content: fmt.Sprintf(` リマインドの期限を更新しました!  id: %d`, todo.ID),
		},
	).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (u *CreateUsecase) Delete(ctx context.Context, todoID int) error {
	// Get the todo info before deleting it
	_, err := u.creator.Get(ctx, todoID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get todo")
		if err.Error() == sql.ErrNoRows.Error() {
			_, _, err = u.traQWSBot.API().MessageApi.PostMessage(ctx, "todo").PostMessageRequest(
				traq.PostMessageRequest{
					Content: fmt.Sprintf(` リマインドが見つかりませんでした!  id: %d`, todoID),
				},
			).Execute()
			if err != nil {
				log.Error().Err(err).Msg("Failed to post message")
				return err
			}
			return nil
		} else {
			return err
		}
	}

	err = u.creator.Delete(ctx, todoID)
	if err != nil {
		return err
	}

	_, _, err = u.traQWSBot.API().MessageApi.PostMessage(ctx, "todo").PostMessageRequest(
		traq.PostMessageRequest{
			Content: fmt.Sprintf(` リマインドを削除しました!  id: %d`, todoID),
		},
	).Execute()
	if err != nil {
		return err
	}

	return nil
}
