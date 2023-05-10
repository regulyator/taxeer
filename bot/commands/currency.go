package commands

import (
	"fmt"
	"log"
	"taxeer/service"
	"time"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type currencyHandler struct {
	command         string
	defaultCurrency string
	bot             *botApi.BotAPI
}

func (handler currencyHandler) CanHandleUpdate(update *botApi.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == handler.command
}

func (handler currencyHandler) HandleUpdate(update *botApi.Update) {
	if handler.CanHandleUpdate(update) {
		currency := handler.defaultCurrency
		var msgText string

		if currentCurrencyRate, err := service.GetCurrencyAtDate(time.Now(), currency); err != nil {
			msgText = fmt.Sprintf("Ooops! Can't get today currency rate for %s, try again later:(", currency)
		} else {
			msgText = fmt.Sprintf("%.4f", currentCurrencyRate)
		}

		if _, err := handler.bot.Send(botApi.NewMessage(update.Message.Chat.ID, msgText)); err != nil {
			log.Println(err)
		}

	}
}

func GetCurrencyCommandHandler(bot *botApi.BotAPI) *currencyHandler {
	return &currencyHandler{
		command:         "currency",
		defaultCurrency: "USD",
		bot:             bot,
	}
}
