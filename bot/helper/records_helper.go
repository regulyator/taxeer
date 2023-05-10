package helper

import (
	"fmt"
	"taxeer/db/sqlc"
	"time"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func WrapRecordsIntoKeyboard(records *[]sqlc.TaxeerRecord, recordDataWrapper func(sqlc.TaxeerRecord, string) string, requestType string) *botApi.InlineKeyboardMarkup {
	var keyboard [][]botApi.InlineKeyboardButton
	for _, record := range *records {
		var wrappedData = recordDataWrapper(record, requestType)
		keyboard = append(keyboard, []botApi.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s %.2f %s", record.Date.Format(time.DateTime), record.IncomeValue, record.IncomeCurrency),
				CallbackData: &wrappedData,
			},
		})
	}
	result := botApi.NewInlineKeyboardMarkup(keyboard...)
	return &result
}
