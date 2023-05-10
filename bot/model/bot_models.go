package model

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateHandler interface {
	CanHandleUpdate(update *botApi.Update) bool
	HandleUpdate(update *botApi.Update)
}

type CallbackRequest struct {
	T string
	D any
}
