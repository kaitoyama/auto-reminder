package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/config"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/db"
	handler "github.com/kaitoyama/kaitoyama-server-template/internal/interface"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func initDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}

func main() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database, err := initDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: cfg.TraqAccessToken,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create bot")
	}

	// Initialize infrastructures
	creator := db.NewTodoCreator(database)
	reminder := db.NewTodoReminder(database)

	// Initialize usecase
	createUsecase := usecase.NewCreateUsecase(creator, bot)
	reminderUsecase := usecase.NewReminderUsecase(bot, reminder)

	// Initialize handler

	handler := handler.NewHandler(createUsecase, reminderUsecase, bot)

	// scheduling goroutine
	go func() {
		s, _ := gocron.NewScheduler()
		_, _ = s.NewJob(
			gocron.DurationJob(
				3*time.Minute,
			),
			gocron.NewTask(
				func() {
					err := reminderUsecase.NotifyTodoInWeek()
					if err != nil {
						log.Error().Err(err).Msg("Failed to notify todo in week")
					}

					err = reminderUsecase.NotifyTodoInThreeDays()
					if err != nil {
						log.Error().Err(err).Msg("Failed to notify todo in three days")
					}

					err = reminderUsecase.NotifyTodoInDay()
					if err != nil {
						log.Error().Err(err).Msg("Failed to notify todo in one day")
					}
				},
			),
		)

		s.Start()
	}()

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		if err := handler.MessageHandler(p); err != nil {
			log.Error().Err(err).Msg("Failed to handle message")
		}
	})

	if err := bot.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

}
