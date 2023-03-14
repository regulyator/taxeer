package service

import (
	"context"
	"database/sql"
	"log"
	"taxeer/db/sqlc"
)

func GetExistUserOrCreate(db *sql.DB, telegramUserId string, chatId int64) *sqlc.TaxeerUser {
	ctx := context.Background()
	query := sqlc.New(db)
	existingUser, err := query.GetUser(ctx, telegramUserId)
	if err == nil {
		if err != nil && existingUser.ChatID != chatId {
			updateUserChatId(existingUser, chatId, query, ctx)
			existingUser.ChatID = chatId
		}
		return &existingUser
	}

	createUserParams := sqlc.CreateUserParams{
		TelegramUserID: telegramUserId,
		ChatID:         chatId,
	}
	createdUser, err := query.CreateUser(ctx, createUserParams)
	if err != nil {
		log.Panic(err)
	}
	return &createdUser
}

func updateUserChatId(existingUser sqlc.TaxeerUser, chatId int64, query *sqlc.Queries, ctx context.Context) {
	updateChatIdParams := sqlc.UpdateUserChatIdParams{
		ID:     existingUser.ID,
		ChatID: chatId,
	}
	if err := query.UpdateUserChatId(ctx, updateChatIdParams); err != nil {
		log.Panic(err)
	}
}
