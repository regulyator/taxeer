package commands

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"taxeer/db/sqlc"
	"taxeer/service"
	"taxeer/util/config"
	"time"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const DDMMYYYYLayout = "02-01-2006"

type incomeHandler struct {
	command  string
	bot      *botApi.BotAPI
	database *sql.DB
}

func (handler incomeHandler) CanHandleUpdate(update *botApi.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == handler.command
}

func (handler incomeHandler) HandleUpdate(update *botApi.Update) {
	if handler.CanHandleUpdate(update) {
		msgText := handleIncomeCommand(update.Message, handler.database)

		if _, err := handler.bot.Send(botApi.NewMessage(update.Message.Chat.ID, msgText)); err != nil {
			log.Println(err)
		}

	}
}

func handleIncomeCommand(message *botApi.Message, database *sql.DB) string {
	if message.Text != "" {
		currentUser := service.GetExistUserOrCreate(database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
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
			Date:           rateDate,
			IncomeValue:    incomeValue,
			IncomeCurrency: incomeParams[1],
			Rate:           currencyRate,
		}
		savedRecord, errCreateRecord := service.CreateIncomeRecord(database, recordParams)
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
		if parsedDate, err := time.Parse(DDMMYYYYLayout, incomeParams[2]); err != nil {
			return time.Now(), err
		} else {
			return parsedDate, nil
		}
	}
}

func GetIncomeCommandHandler(bot *botApi.BotAPI, postgresDb *config.PostgresDb) *incomeHandler {
	return &incomeHandler{
		command:  "income",
		bot:      bot,
		database: postgresDb.Database,
	}
}
