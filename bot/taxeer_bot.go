package bot

import (
	"log"
	"os"
	"taxeer/bot/commands"
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
	bot.Debug = true
	updateConfig := botApi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}
		var msgText = handleCommand(update.Message, postgresDb)
		if len(msgText) > 0 {
			msg := botApi.NewMessage(update.Message.Chat.ID, msgText)
			sendResponse(msg, bot)
		}
	}
}

func registerCommands(bot *botApi.BotAPI) {
	commandsConfig := botApi.NewSetMyCommands(
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
		},
		botApi.BotCommand{
			Command:     "/current",
			Description: "Print current incomes sum year, mont and month taxes values",
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

func handleCommand(message *botApi.Message, postgresDb *config.PostgresDb) string {
	command := message.Command()
	switch command {
	case "income":
		return commands.HandleIncomeCommand(message, postgresDb)
	case "currency":
		return commands.HandleCurrencyCommand("USD")
	case "statistic":
		return commands.HandleStatisticCommand(message, postgresDb)
	case "start":
		return "Hi! I'm bot for Georgian accounting, feel free to see commands list for usage!"
	case "current":
		return commands.HandleCurrentCommand(message, postgresDb)
	default:
		return ""
	}
}
