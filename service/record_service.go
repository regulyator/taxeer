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

func GetAllUSerRecordsInCurrentYear(db *sql.DB, telegramUserId string, chatId int64) (*[]sqlc.TaxeerRecord, error) {
	currentDate := time.Now()
	return getUserRecordsByDateBetween(
		db,
		telegramUserId,
		chatId,
		time.Date(currentDate.Year(), time.January, 1, 0, 0, 0, 0, time.Local),
		currentDate)
}

func GetAllUSerRecordsInCurrentMonth(db *sql.DB, telegramUserId string, chatId int64) (*[]sqlc.TaxeerRecord, error) {
	currentDate := time.Now()
	return getUserRecordsByDateBetween(
		db,
		telegramUserId,
		chatId,
		time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, time.Local),
		currentDate)
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
