package commands

import (
	"fmt"
	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"taxeer/service"
	"taxeer/util/config"
	"time"
)

func HandleStatisticCommand(message *botApi.Message, postgresDb *config.PostgresDb) string {
	records, err := service.GetLastTenUserRecords(postgresDb.Database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
	if err != nil || len(*records) == 0 {
		return "Sorry can't retrieve records:( Maybe you not save any income yet?"
	}
	var result []string
	for _, record := range *records {
		result = append(result, fmt.Sprintf("%s %.2f %s", record.Date.Format(time.DateTime), record.IncomeValue, record.IncomeCurrency))
	}
	return strings.Join(result, "\n")
}
