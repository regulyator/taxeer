package commands

import (
	"fmt"
	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"taxeer/db/sqlc"
	"taxeer/service"
	"taxeer/util/config"
)

const MonthTaxPercentage = 0.01

type Incomes struct {
	yearSum  float64
	monthSum float64
}

func HandleCurrentCommand(message *botApi.Message, postgresDb *config.PostgresDb) string {
	yearRecords, err := service.GetAllUSerRecordsInCurrentFinanceYear(postgresDb.Database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
	monthRecords, err := service.GetAllUSerRecordsInCurrentFinanceMonth(postgresDb.Database, strconv.FormatInt(message.From.ID, 10), message.Chat.ID)
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
