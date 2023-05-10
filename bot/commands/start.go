package commands

import (
	"log"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type startHandler struct {
	command string
	bot     *botApi.BotAPI
}

func (handler startHandler) CanHandleUpdate(update *botApi.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == handler.command
}

func (handler startHandler) HandleUpdate(update *botApi.Update) {
	if handler.CanHandleUpdate(update) {

		if _, err := handler.bot.Send(botApi.NewMessage(update.Message.Chat.ID, "Hi! I'm bot for Georgian accounting, feel free to see commands list for usage!")); err != nil {
			log.Println(err)
		}

	}
}

func GetStartCommandHandler(bot *botApi.BotAPI) *startHandler {
	return &startHandler{
		command: "start",
		bot:     bot,
	}
}
