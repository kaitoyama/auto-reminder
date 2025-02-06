package usecase

import (
	"context"
	"fmt"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

type ReminderUsecase struct {
	traqWSBot *traqwsbot.Bot
	reminder  domain.TodoReminder
}

func NewReminderUsecase(traqWSBot *traqwsbot.Bot, reminder domain.TodoReminder) *ReminderUsecase {
	return &ReminderUsecase{
		traqWSBot: traqWSBot,
		reminder:  reminder,
	}
}

func (u *ReminderUsecase) NotifyTodoInWeek() error {
	ctx := context.Background()
	todos, err := u.reminder.GetTodoInWeek(ctx)
	if err != nil {
		return err
	}

	for _, todo := range todos {
		// Notify todo to the channel
		users := ""
		for _, user := range todo.Users {
			users += fmt.Sprintf("@%s ", user)
		}
		_, _, err = u.traqWSBot.API().MessageApi.PostMessage(ctx, todo.ChannelID).PostMessageRequest(
			traq.PostMessageRequest{
				Content: fmt.Sprintf(` ### 一週間前のリマインドです!
%s cc: @%s
%s
期限: %s
`, users, todo.OwnerID, todo.Content, todo.DueAt.Format("2006-01-02")),
			},
		).Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *ReminderUsecase) NotifyTodoInThreeDays() error {
	ctx := context.Background()
	todos, err := u.reminder.GetTodoInThreeDays(ctx)
	if err != nil {
		return err
	}

	for _, todo := range todos {
		// Notify todo to the channel
		users := ""
		for _, user := range todo.Users {
			users += fmt.Sprintf("@%s ", user)
		}
		_, _, err = u.traqWSBot.API().MessageApi.PostMessage(ctx, todo.ChannelID).PostMessageRequest(
			traq.PostMessageRequest{
				Content: fmt.Sprintf(` ### 三日前のリマインドです!
%s cc: @%s
%s
期限: %s
`, users, todo.OwnerID, todo.Content, todo.DueAt.Format("2006-01-02")),
			},
		).Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *ReminderUsecase) NotifyTodoInDay() error {
	ctx := context.Background()
	todos, err := u.reminder.GetTodoInDay(ctx)
	if err != nil {
		return err
	}

	for _, todo := range todos {
		// Notify todo to the channel
		users := ""
		for _, user := range todo.Users {
			users += fmt.Sprintf("@%s ", user)
		}
		_, _, err = u.traqWSBot.API().MessageApi.PostMessage(ctx, todo.ChannelID).PostMessageRequest(
			traq.PostMessageRequest{
				Content: fmt.Sprintf(` ### 本日のリマインドです!
%s cc: @%s
%s
期限: %s
`, users, todo.OwnerID, todo.Content, todo.DueAt.Format("2006-01-02")),
			},
		).Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
