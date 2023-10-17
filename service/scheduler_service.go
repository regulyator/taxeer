package service

import (
	"context"
	"database/sql"
	"github.com/go-co-op/gocron"
	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"taxeer/db/sqlc"
	"time"
)

var goCronScheduler *gocron.Scheduler

func ScheduleTaxNotification(bot *botApi.BotAPI, database *sql.DB) {
	if goCronScheduler == nil {
		goCronScheduler = gocron.NewScheduler(time.UTC)
	}

	createScheduledNotificationJob(
		"0 0 7 * *",
		func() {
			notifyAllUsersAboutTaxPeriod(database, bot, "Please, dont forget to fill out tax monthly income due the 15th of this month!")
		})

	createScheduledNotificationJob(
		"0 0 14 * *",
		func() {
			notifyAllUsersAboutTaxPeriod(database, bot, "Please, dont forget to fill out tax monthly income due the 15th of this month!")
		})

	createScheduledNotificationJob(
		"0 0 L-2 * *",
		func() {
			notifyAllUsersAboutTaxPeriod(database, bot, "Please, dont forget to fill out monthly income of this month!")
		})

	goCronScheduler.StartAsync()
}

func createScheduledNotificationJob(cronPattern string, task func()) {
	if _, err := goCronScheduler.Cron("").Do(task); err != nil {
		log.Printf("error when schedulling task %s cron: %s", err, cronPattern)
	}
}

func notifyAllUsersAboutTaxPeriod(db *sql.DB, bot *botApi.BotAPI, msgText string) {
	ctx := context.Background()
	query := sqlc.New(db)
	users, err := query.GetAllUsers(ctx)

	if err != nil {
		log.Printf("error when retrieving all users %s", err)
		return
	}
	for _, u := range users {
		if _, err := bot.Send(botApi.NewMessage(u.ChatID, msgText)); err != nil {
			log.Printf("error when sending message to user: %s, message: %s, error: %s", u.TelegramUserID, msgText, err)
		}
	}
}
