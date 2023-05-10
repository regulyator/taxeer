package commands

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"taxeer/db/sqlc"
	"taxeer/service"
	"taxeer/util/config"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const MonthTaxPercentage = 0.01

type Incomes struct {
	yearSum  float64
	monthSum float64
}

type currentHandler struct {
	command  string
	bot      *botApi.BotAPI
	database *sql.DB
}

func (handler currentHandler) CanHandleUpdate(update *botApi.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == handler.command
}

func (handler currentHandler) HandleUpdate(update *botApi.Update) {
	if handler.CanHandleUpdate(update) {
		msgText := handleCurrentCommand(update.Message, handler.database)

		if _, err := handler.bot.Send(botApi.NewMessage(update.Message.Chat.ID, msgText)); err != nil {
			log.Println(err)
		}

	}
}

func handleCurrentCommand(message *botApi.Message, database *sql.DB) string {
	yearRecords, err := service.GetAllUSerRecordsInCurrentFinanceYear(database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
	monthRecords, err := service.GetAllUSerRecordsInCurrentFinanceMonth(database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
	if err != nil || len(*yearRecords) == 0 {
		return "Looks like you no need to pay taxes in this month:) Check saved incomes by /statistic command!"
	}
	incomes := calculateYearAndMonthIncomes(yearRecords, monthRecords)

	return fmt.Sprintf("This year incomes sum: %.2f\nThis month incomes sum: %.2f\nThis month taxes: %.2f\nAll values in GEL!", incomes.yearSum, incomes.monthSum, incomes.monthSum*MonthTaxPercentage)
}

func calculateYearAndMonthIncomes(yearRecords *[]sqlc.TaxeerRecord, monthRecords *[]sqlc.TaxeerRecord) Incomes {
	var yearSum float64
	for _, record := range *yearRecords {
		yearSum += record.Rate * record.IncomeValue
	}

	var monthSum float64
	for _, record := range *monthRecords {
		monthSum += record.Rate * record.IncomeValue
	}

	return Incomes{
		yearSum:  yearSum,
		monthSum: monthSum,
	}
}

func GetCurrentCommandHandler(bot *botApi.BotAPI, postgresDb *config.PostgresDb) *currentHandler {
	return &currentHandler{
		command:  "current",
		bot:      bot,
		database: postgresDb.Database,
	}
}
