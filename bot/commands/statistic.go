package commands

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"taxeer/service"
	"taxeer/util/config"
	"time"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type statisticHandler struct {
	command  string
	bot      *botApi.BotAPI
	database *sql.DB
}

func (handler statisticHandler) CanHandleUpdate(update *botApi.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == handler.command
}

func (handler statisticHandler) HandleUpdate(update *botApi.Update) {
	if handler.CanHandleUpdate(update) {
		var msgText string

		if records, err := service.GetLastTenUserRecords(handler.database, strconv.FormatInt(update.Message.From.ID, 10), update.Message.Chat.ID); err != nil || len(*records) == 0 {
			msgText = "Sorry can't retrieve records:( Maybe you not save any income yet?"
		} else {
			var result []string
			for _, record := range *records {
				result = append(result, fmt.Sprintf("%s %.2f %s", record.Date.Format(time.DateTime), record.IncomeValue, record.IncomeCurrency))
			}
			msgText = strings.Join(result, "\n")
		}

		if _, err := handler.bot.Send(botApi.NewMessage(update.Message.Chat.ID, msgText)); err != nil {
			log.Println(err)
		}

	}
}

func GetStatisticCommandHandler(bot *botApi.BotAPI, postgresDb *config.PostgresDb) *statisticHandler {
	return &statisticHandler{
		command:  "statistic",
		bot:      bot,
		database: postgresDb.Database,
	}
}
