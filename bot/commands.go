package bot

import (
	"fmt"
	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"taxeer/db/sqlc"
	"taxeer/service"
	"taxeer/util/config"
	"time"
)

const layoutDateTime = "2006-01-02 15:04"

func HandleOtherCommand() string {
	return "Sorry, i can't help with this:("
}

func HandleCurrencyCommand(currency string) string {
	currentCurrencyRate, err := service.GetCurrencyAtDate(time.Now(), currency)
	if err != nil {
		return fmt.Sprintf("Ooops! Can't get today currency rate for %s, try again later:(", currency)
	}
	return fmt.Sprintf("%.4f", currentCurrencyRate)
}

func HandleIncomeCommand(message *botApi.Message, postgresDb *config.PostgresDb) string {
	if message.Text != "" {
		currentUser := service.GetExistUserOrCreate(postgresDb.Database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
		incomeParams := strings.Split(strings.TrimSpace(strings.Replace(message.Text, "/income", "", 1)), ":")
		incomeValue, err := strconv.ParseFloat(incomeParams[0], 64)
		if err != nil || len(incomeParams) != 2 || len(incomeParams[1]) == 0 {
			return "Please input correct income value after command:)"
		}
		currencyRate, errGettingRate := service.GetCurrencyAtDate(time.Now(), incomeParams[1])
		if errGettingRate != nil {
			return fmt.Sprintf("Error getting currency rate for %s Income not saved!:( Please try again later...", incomeParams[1])
		}
		recordParams := sqlc.CreateRecordParams{
			TaxeerUserID:   currentUser.ID,
			Date:           time.Now(),
			IncomeValue:    incomeValue,
			IncomeCurrency: incomeParams[1],
			Rate:           currencyRate,
		}
		savedRecord, errCreateRecord := service.CreateIncomeRecord(postgresDb.Database, recordParams)
		if errCreateRecord != nil {
			return "Income not saved, try again:("
		}
		return fmt.Sprintf("Income %.2f in %s saved!", savedRecord.IncomeValue, savedRecord.IncomeCurrency)
	} else {
		return "Please input correct income value after command:)"
	}
}

func HandleStatisticCommand(message *botApi.Message, postgresDb *config.PostgresDb) string {
	records, err := service.GetLastTenUserRecords(postgresDb.Database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
	if err != nil || len(*records) == 0 {
		return `Sorry can't retrieve records:(
				Maybe you not save any income yet?`
	}
	var result []string
	for _, record := range *records {
		result = append(result, fmt.Sprintf("%s %.2f %s", record.Date.Format(layoutDateTime), record.IncomeValue, record.IncomeCurrency))
	}
	return strings.Join(result, "\n")
}
