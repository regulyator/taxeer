package bot

import (
	"log"
	"os"
	"taxeer/bot/commands"
	"taxeer/bot/model"
	"taxeer/util/config"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartTaxeerBotListener(postgresDb *config.PostgresDb) {
	bot, err := botApi.NewBotAPI(os.Getenv("BOT_API_KEY"))
	if err != nil {
		log.Printf("Error when starting bot!")
		log.Panic(err)
	}
	registerCommands(bot)
	handleUpdates(postgresDb, bot)
}

func handleUpdates(postgresDb *config.PostgresDb, bot *botApi.BotAPI) {
	handlers := getHandlers(postgresDb, bot)
	bot.Debug = true
	updateConfig := botApi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		for _, handler := range handlers {
			if handler.CanHandleUpdate(&update) {
				handler.HandleUpdate(&update)
			}
		}

	}
}

func getHandlers(postgresDb *config.PostgresDb, bot *botApi.BotAPI) []model.UpdateHandler {
	handlers := []model.UpdateHandler{
		commands.GetCurrencyCommandHandler(bot),
		commands.GetCurrentCommandHandler(bot, postgresDb),
		commands.GetDeleteCommandHandler(bot, postgresDb),
		commands.GetIncomeCommandHandler(bot, postgresDb),
		commands.GetStartCommandHandler(bot),
		commands.GetStatisticCommandHandler(bot, postgresDb),
	}
	return handlers
}

func registerCommands(bot *botApi.BotAPI) {
	commandsConfig := botApi.NewSetMyCommands(
		botApi.BotCommand{
			Command:     "/income",
			Description: "Save income by command '/income value:currency:date'\n(don't forget separate by ':')\ndate params is optional (DD-MM-YYYY)",
		},
		botApi.BotCommand{
			Command:     "/currency",
			Description: "Today currency rate for USD",
		},
		botApi.BotCommand{
			Command:     "/statistic",
			Description: "Print last 10 incomes",
		},
		botApi.BotCommand{
			Command:     "/current",
			Description: "Print current incomes sum year, mont and month taxes values",
		},
		botApi.BotCommand{
			Command:     "/delete",
			Description: "Delete selected income",
		})
	if _, err := bot.Request(commandsConfig); err != nil {
		log.Panic(err)
	}
}

func sendResponse(msg botApi.MessageConfig, bot *botApi.BotAPI) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
