package service

import (
	"context"
	"database/sql"
	"log"
	"taxeer/db/sqlc"
	"time"
)

func CreateIncomeRecord(db *sql.DB, recordParams sqlc.CreateRecordParams) (*sqlc.TaxeerRecord, error) {
	ctx := context.Background()
	query := sqlc.New(db)
	createdRecord, err := query.CreateRecord(ctx, recordParams)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return &createdRecord, nil
}

func GetLastTenUserRecords(db *sql.DB, telegramUserId string, chatId int64) (*[]sqlc.TaxeerRecord, error) {
	ctx := context.Background()
	query := sqlc.New(db)
	currentUser := GetExistUserOrCreate(db, telegramUserId, chatId)
	requestParams := sqlc.GetLastNRecordByUserIdParams{
		TaxeerUserID: currentUser.ID,
		Limit:        10,
	}
	records, err := query.GetLastNRecordByUserId(ctx, requestParams)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &records, nil
}

func GetAllUSerRecordsInCurrentFinanceYear(db *sql.DB, telegramUserId string, chatId int64) (*[]sqlc.TaxeerRecord, error) {
	currentDate := time.Now()
	currentYear, currentMonth, _ := currentDate.Date()
	var dateFrom, dateTo time.Time
	if currentMonth == time.January {
		dateFrom = time.Date(currentYear-1, time.January, 1, 0, 0, 0, 0, time.Local)
		dateTo = time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.Local)
	} else {
		dateFrom = time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.Local)
		dateTo = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
	}
	return getUserRecordsByDateBetween(
		db,
		telegramUserId,
		chatId,
		dateFrom,
		dateTo)
}

func GetAllUSerRecordsInCurrentFinanceMonth(db *sql.DB, telegramUserId string, chatId int64) (*[]sqlc.TaxeerRecord, error) {
	currentDate := time.Now()
	currentYear, currentMonth, _ := currentDate.Date()
	var dateFrom, dateTo time.Time
	if currentMonth == time.January {
		dateFrom = time.Date(currentYear-1, time.December, 1, 0, 0, 0, 0, time.Local)
		dateTo = time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.Local)
	} else {
		dateFrom = time.Date(currentYear, currentMonth+time.Month(-1), 1, 0, 0, 0, 0, time.Local)
		dateTo = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
	}
	return getUserRecordsByDateBetween(
		db,
		telegramUserId,
		chatId,
		dateFrom,
		dateTo)
}

func getUserRecordsByDateBetween(db *sql.DB, telegramUserId string, chatId int64, dateFrom time.Time, dateTo time.Time) (*[]sqlc.TaxeerRecord, error) {
	ctx := context.Background()
	query := sqlc.New(db)
	currentUser := GetExistUserOrCreate(db, telegramUserId, chatId)
	requestParams := sqlc.GetRecordByUserIdAndDateBetweenOrderedByDateDescParams{
		TaxeerUserID: currentUser.ID,
		Date:         dateFrom,
		Date_2:       dateTo,
	}
	records, err := query.GetRecordByUserIdAndDateBetweenOrderedByDateDesc(ctx, requestParams)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &records, nil
}
