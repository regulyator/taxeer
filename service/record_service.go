package service

import (
	"context"
	"database/sql"
	"log"
	"taxeer/db/sqlc"
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
