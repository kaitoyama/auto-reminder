package handler

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/config"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/rs/zerolog/log"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

type Handler struct {
	create   *usecase.CreateUsecase
	reminder *usecase.ReminderUsecase
	bot      *traqwsbot.Bot
}

func (h *Handler) MessageHandler(p *payload.MessageCreated) error {

	// configからbotの名前を取得
	config := config.LoadConfig()

	if !strings.HasPrefix(p.Message.PlainText, "@"+config.BotName) {
		log.Info().Msg("Message received")
		return nil
	}

	// spaceで区切る
	args := strings.Fields(p.Message.PlainText)
	if len(args) < 2 {
		_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
			traq.PostMessageRequest{
				Content: "コマンドが不正です",
			},
		).Execute()
		if err != nil {
			log.Error().Err(err).Msg("Failed to post message")
			return err
		}
	}

	switch args[1] {
	case "create":
		// createコマンドの処理
		// -u @user1 @user2 -d 2021-01-01 -c content
		// という形式で受け取る
		// それぞれのオプションは順不同で受け取る
		// -u: ユーザーを指定 (複数指定可能)
		// -d: 期限を指定
		// -c: 内容を指定
		// それぞれのオプションは必須

		var users []string
		var deadline string
		var dueAt time.Time
		var content string

		i := 2
		for i < len(args) {
			switch args[i] {
			case "-u":
				i++
				for i < len(args) && !strings.HasPrefix(args[i], "-") {
					users = append(users, strings.Split(args[i], "@")[1])
					i++
				}
			case "-d":
				i++
				if i >= len(args) {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				deadline = args[i]
				parsed, err := time.Parse("2006-01-02", deadline)
				dueAt = parsed
				if err != nil {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				i++
			case "-c":
				i++
				if i >= len(args) {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				// Join remaining tokens as content.
				content = strings.Join(args[i:], " ")
				i = len(args)
			default:
				i++
			}
		}

		// Validate required options.
		if len(users) == 0 || deadline == "" || content == "" {
			_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "コマンドが不正です",
				},
			).Execute()
			if err != nil {
				log.Error().Err(err).Msg("Failed to post message")
				return err
			}
			return nil
		}

		// Process the create command using the create usecase.
		if err := h.create.Create(context.Background(), p.Message.ChannelID, content, dueAt, users, p.Message.User.Name); err != nil {
			_, _, msgErr := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "処理中にエラーが発生しました",
				},
			).Execute()
			if msgErr != nil {
				log.Error().Err(msgErr).Msg("Failed to post message")
				return msgErr
			}
			log.Error().Err(err).Msg("Failed to create todo")
			return err
		}

	case "addUser":
		// addUserコマンドの処理
		// -u @user1 @user2 -i id
		// という形式で受け取る
		// それぞれのオプションは順不同で受け取る
		// -u: ユーザーを指定 (複数指定可能)
		// -i: IDを指定
		// それぞれのオプションは必須

		var users []string
		var id int

		i := 2
		for i < len(args) {
			switch args[i] {
			case "-u":
				i++
				for i < len(args) && !strings.HasPrefix(args[i], "-") {
					users = append(users, strings.Split(args[i], "@")[1])
					i++
				}
			case "-i":
				i++
				if i >= len(args) {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				parsed, err := strconv.Atoi(args[i])
				if err != nil {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				id = parsed
				i++
			default:
				i++
			}
		}

		// Validate required options.
		if len(users) == 0 || id == 0 {
			_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "コマンドが不正です",
				},
			).Execute()
			if err != nil {
				log.Error().Err(err).Msg("Failed to post message")
				return err
			}
			return nil
		}

		// Process the addUser command using the create usecase.
		if err := h.create.AddUser(context.Background(), id, users); err != nil {
			_, _, msgErr := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "処理中にエラーが発生しました",
				},
			).Execute()
			if msgErr != nil {
				log.Error().Err(msgErr).Msg("Failed to post message")
				return msgErr
			}
			log.Error().Err(err).Msg("Failed to add user")
			return err
		}

	case "deleteUser":
		// deleteUserコマンドの処理
		// -u @user -i id
		// 一度に削除できるのは1人のみです

		var users []string
		var id int

		i := 2
		for i < len(args) {
			switch args[i] {
			case "-u":
				i++
				if i >= len(args) {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				// collect user tokens
				for i < len(args) && !strings.HasPrefix(args[i], "-") {
					users = append(users, strings.Split(args[i], "@")[1])
					i++
				}
			case "-i":
				i++
				if i >= len(args) {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				parsed, err := strconv.Atoi(args[i])
				if err != nil {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				id = parsed
				i++
			default:
				i++
			}
		}

		// ユーザーは1人のみ許可
		if len(users) != 1 || id == 0 {
			_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "コマンドが不正です",
				},
			).Execute()
			if err != nil {
				log.Error().Err(err).Msg("Failed to post message")
				return err
			}
			return nil
		}

		user := users[0]

		// Process the deleteUser command using the create usecase.
		if err := h.create.DeleteUser(context.Background(), id, user); err != nil {
			_, _, msgErr := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "処理中にエラーが発生しました",
				},
			).Execute()
			if msgErr != nil {
				log.Error().Err(msgErr).Msg("Failed to post message")
				return msgErr
			}
			log.Error().Err(err).Msg("Failed to delete user")
			return err
		}

	case "updateDueAt":
		// updateDueAtコマンドの処理
		// -d 2021-01-01 -i id
		// という形式で受け取る
		// それぞれのオプションは順不同で受け取る
		// -d: 期限を指定
		// -i: IDを指定
		// それぞれのオプションは必須

		var deadline string
		var dueAt time.Time
		var id int

		i := 2
		for i < len(args) {
			switch args[i] {
			case "-d":
				i++
				if i >= len(args) {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				deadline = args[i]
				parsed, err := time.Parse("2006-01-02", deadline)
				dueAt = parsed
				if err != nil {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				i++
			case "-i":
				i++
				if i >= len(args) {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				parsed, err := strconv.Atoi(args[i])
				if err != nil {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				id = parsed
				i++
			default:
				i++
			}
		}

		// Validate required options.
		if deadline == "" || id == 0 {
			_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "コマンドが不正です",
				},
			).Execute()
			if err != nil {
				log.Error().Err(err).Msg("Failed to post message")
				return err
			}
			return nil
		}

		// Process the updateDueAt command using the create usecase.
		if err := h.create.UpdateDueAt(context.Background(), id, dueAt); err != nil {
			_, _, msgErr := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "処理中にエラーが発生しました",
				},
			).Execute()
			if msgErr != nil {
				log.Error().Err(msgErr).Msg("Failed to post message")
				return msgErr
			}
			log.Error().Err(err).Msg("Failed to update due date")
			return err
		}

	case "delete":
		// deleteコマンドの処理
		// -i id
		// という形式で受け取る
		// すべてのオプションは必須

		var id int

		i := 2
		for i < len(args) {
			switch args[i] {
			case "-i":
				i++
				if i >= len(args) {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				parsed, err := strconv.Atoi(args[i])
				if err != nil {
					_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
						traq.PostMessageRequest{
							Content: "コマンドが不正です",
						},
					).Execute()
					if err != nil {
						log.Error().Err(err).Msg("Failed to post message")
						return err
					}
					return nil
				}
				id = parsed
				i++
			default:
				i++
			}
		}

		// Validate required options.
		if id == 0 {
			_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "コマンドが不正です",
				},
			).Execute()
			if err != nil {
				log.Error().Err(err).Msg("Failed to post message")
				return err
			}
			return nil
		}

		// Process the delete command using the create usecase.
		if err := h.create.Delete(context.Background(), id); err != nil {
			_, _, msgErr := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
				traq.PostMessageRequest{
					Content: "処理中にエラーが発生しました",
				},
			).Execute()
			if msgErr != nil {
				log.Error().Err(msgErr).Msg("Failed to post message")
				return msgErr
			}
			log.Error().Err(err).Msg("Failed to delete todo")
			return err
		}

	default:
		_, _, err := h.bot.API().MessageApi.PostMessage(context.Background(), p.Message.ChannelID).PostMessageRequest(
			traq.PostMessageRequest{
				Content: "コマンドが不正です",
			},
		).Execute()
		if err != nil {
			log.Error().Err(err).Msg("Failed to post message")
			return err
		}
	}

	return nil
}

func NewHandler(create *usecase.CreateUsecase, reminder *usecase.ReminderUsecase, bot *traqwsbot.Bot) *Handler {
	return &Handler{
		create:   create,
		reminder: reminder,
		bot:      bot,
	}
}
