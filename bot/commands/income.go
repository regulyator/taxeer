package commands

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

func HandleIncomeCommand(message *botApi.Message, postgresDb *config.PostgresDb) string {
	if message.Text != "" {
		currentUser := service.GetExistUserOrCreate(postgresDb.Database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
		incomeParams := strings.Split(strings.TrimSpace(strings.Replace(message.Text, "/income", "", 1)), ":")
		incomeValue, err := strconv.ParseFloat(incomeParams[0], 64)
		if err != nil || len(incomeParams) < 2 || len(incomeParams[1]) == 0 {
			return "Please input correct income value after command:)"
		}

		rateDate, errParseDate := parseDateInput(incomeParams)
		if errParseDate != nil {
			return "Ooops, some error! If you use date in income parameters it should be in 'DD-MM-YYYY' format."
		}

		currencyRate, errGettingRate := service.GetCurrencyAtDate(rateDate, incomeParams[1])
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

func parseDateInput(incomeParams []string) (time.Time, error) {
	if len(incomeParams) == 2 {
		return time.Now(), nil
	} else {
		if parsedDate, err := time.Parse(time.DateOnly, incomeParams[2]); err != nil {
			return time.Now(), err
		} else {
			return parsedDate, nil
		}
	}
}
