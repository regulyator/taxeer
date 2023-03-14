package bot

import (
	"log"
	"os"
	"taxeer/util/config"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartTaxeerBotListener(postgresDb *config.PostgresDb) {
	println(os.Getenv("BOT_API_KEY"))
	bot, err := botApi.NewBotAPI(os.Getenv("BOT_API_KEY"))
	if err != nil {
		log.Panic(err)
	}
	registerCommands(bot)
	handleUpdates(postgresDb, bot)
}

func handleUpdates(postgresDb *config.PostgresDb, bot *botApi.BotAPI) {
	bot.Debug = true
	updateConfig := botApi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := botApi.NewMessage(update.Message.Chat.ID, handleCommand(update.Message, postgresDb))
		sendResponse(msg, bot)
	}
}

func registerCommands(bot *botApi.BotAPI) {
	commands := botApi.NewSetMyCommands(
		botApi.BotCommand{
			Command:     "/income",
			Description: "Save income by command '/income value:currency' (don't forget separate by ':')",
		},
		botApi.BotCommand{
			Command:     "/currency",
			Description: "Today currency rate for USD",
		},
		botApi.BotCommand{
			Command:     "/statistic",
			Description: "Print last 10 incomes",
		})
	if _, err := bot.Request(commands); err != nil {
		log.Panic(err)
	}
}

func sendResponse(msg botApi.MessageConfig, bot *botApi.BotAPI) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func handleCommand(message *botApi.Message, postgresDb *config.PostgresDb) string {
	command := message.Command()
	switch command {
	case "income":
		return HandleIncomeCommand(message, postgresDb)
	case "currency":
		return HandleCurrencyCommand("USD")
	case "statistic":
		return HandleStatisticCommand(message, postgresDb)
	default:
		return HandleOtherCommand()
	}
}
