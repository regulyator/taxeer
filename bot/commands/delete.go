package commands

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"taxeer/bot/helper"
	"taxeer/bot/model"
	"taxeer/db/sqlc"
	"taxeer/service"
	"taxeer/util/config"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

const (
	DELETE_PAGE    string = "dp"
	DELETE_COMMAND string = "del"
)

type PageRequest struct {
	Lim int32
	Off int32
}

type DeleteRequest struct {
	ID uuid.UUID
}

type deleteHandler struct {
	command  string
	bot      *botApi.BotAPI
	database *sql.DB
}

func (handler deleteHandler) CanHandleUpdate(update *botApi.Update) bool {
	return (update.Message != nil && update.Message.IsCommand() && update.Message.Command() == handler.command) || isDeleteCallback(update)
}

func (handler deleteHandler) HandleUpdate(update *botApi.Update) {
	if handler.CanHandleUpdate(update) {
		if update.Message != nil && update.Message.IsCommand() {
			log.Println("AAAAAAAAAAAAAAAAAA")
			if deletionResponseMarkup, err := getMessageForDeletionRecords(update.Message.From.ID, update.Message.Chat.ID, handler.database, 10, 0, DELETE_COMMAND); err == nil {
				response := botApi.NewMessage(update.Message.Chat.ID, "Select record for deletion")
				response.ReplyMarkup = deletionResponseMarkup

				if _, err := handler.bot.Send(response); err != nil {
					log.Println(err)
				}
			}
		} else if isDeleteCallback(update) {
			log.Println("BBBBBBBBBBBBBBBBBBB")
			processDeleteCallback(handler, update)
		}

	}
}

func processDeleteCallback(handler deleteHandler, update *botApi.Update) {
	var callbackRequest model.CallbackRequest
	if unmarshalError := json.Unmarshal([]byte(update.CallbackQuery.Data), &callbackRequest); unmarshalError == nil {
		switch requestType := callbackRequest.T; requestType {
		case DELETE_PAGE:
			log.Println("HANDLE DELETE_PAGE")
			handlePageRequest(callbackRequest, update, handler)
		case DELETE_COMMAND:
			log.Println("HANDLE DELETE_COMMAND")
			handleDeleteRequest(callbackRequest, update, handler)
		default:
			callbackErrorResponse(update.CallbackQuery.ID, handler.bot)
		}
	} else {
		callbackErrorResponse(update.CallbackQuery.ID, handler.bot)
	}
}

func handlePageRequest(callbackRequest model.CallbackRequest, update *botApi.Update, handler deleteHandler) {
	if pageRequest, err := extractPayload[PageRequest](callbackRequest.D); err != nil {
		callbackErrorResponse(update.CallbackQuery.ID, handler.bot)
	} else {
		if deletionResponseMarkup, err := getMessageForDeletionRecords(update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, handler.database, pageRequest.Lim, pageRequest.Off, DELETE_COMMAND); err == nil {
			response := botApi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, *deletionResponseMarkup)

			var callbackResponse botApi.CallbackConfig
			if _, err := handler.bot.Request(response); err != nil {
				log.Println(err)
				callbackResponse = botApi.NewCallback(update.CallbackQuery.ID, "Can't process callback!")
			} else {
				callbackResponse = botApi.NewCallback(update.CallbackQuery.ID, "Records loaded...")
			}

			handler.bot.Request(callbackResponse)
		}
	}
}

func handleDeleteRequest(callbackRequest model.CallbackRequest, update *botApi.Update, handler deleteHandler) {
	if deleteRequest, err := extractPayload[DeleteRequest](callbackRequest.D); err != nil {
		callbackErrorResponse(update.CallbackQuery.ID, handler.bot)
	} else {
		if errDelete := service.DeleteUserRecord(handler.database, deleteRequest.ID); errDelete == nil {
			handler.bot.Request(botApi.NewCallback(update.CallbackQuery.ID, "Record deleted!"))
			response := botApi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID)
			handler.bot.Request(response)
		} else {
			log.Println(errDelete)
			handler.bot.Request(botApi.NewCallback(update.CallbackQuery.ID, "Can't delete record!"))
			response := botApi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID)
			handler.bot.Request(response)
		}
	}
}

func extractPayload[R any](payload any) (*R, error) {
	if marshaledData, err := json.Marshal(payload); err != nil {
		return nil, err
	} else {
		var resultData R
		if err := json.Unmarshal(marshaledData, &resultData); err != nil {
			return nil, err
		} else {
			return &resultData, nil
		}
	}
}

func getMessageForDeletionRecords(userID int64, chatID int64, database *sql.DB, limit int32, offset int32, requestType string) (*botApi.InlineKeyboardMarkup, error) {
	records, err := service.GetUserRecordsPage(database, strconv.FormatInt(userID, 10), chatID, limit, offset)

	if err != nil || len(*records) == 0 {
		return nil, err
	}

	replyMarkup := helper.WrapRecordsIntoKeyboard(records, wrapRecordForDeleteRequest, requestType)
	addNavigationButton(limit, offset, replyMarkup)

	return replyMarkup, nil
}

func addNavigationButton(limit int32, offset int32, replyMarkup *botApi.InlineKeyboardMarkup) {
	prevPageRequest, _ := json.Marshal(model.CallbackRequest{
		T: DELETE_PAGE,
		D: PageRequest{limit, offset - limit},
	})
	nextPageRequest, _ := json.Marshal(model.CallbackRequest{
		T: DELETE_PAGE,
		D: PageRequest{limit, offset + limit},
	})

	replyMarkup.InlineKeyboard = append(replyMarkup.InlineKeyboard, botApi.NewInlineKeyboardRow(
		botApi.NewInlineKeyboardButtonData("Back", string(prevPageRequest)),
		botApi.NewInlineKeyboardButtonData("Next", string(nextPageRequest)),
	))
}

func isDeleteCallback(update *botApi.Update) bool {
	if update.CallbackQuery == nil {
		return false
	}

	var callbackRequest model.CallbackRequest
	if unmarshalError := json.Unmarshal([]byte(update.CallbackQuery.Data), &callbackRequest); unmarshalError == nil {
		return callbackRequest.T == DELETE_PAGE || callbackRequest.T == DELETE_COMMAND
	} else {
		return false
	}

}

func wrapRecordForDeleteRequest(record sqlc.TaxeerRecord, requestType string) string {
	data, _ := json.Marshal(model.CallbackRequest{
		T: requestType,
		D: DeleteRequest{record.ID},
	})
	println("DATA LENGTH")
	println(len(data))
	return string(data)
}

func callbackErrorResponse(callbackID string, botAPI *botApi.BotAPI) {
	callbackResponse := botApi.NewCallback(callbackID, "Can't process callback!")
	botAPI.Request(callbackResponse)
}

func GetDeleteCommandHandler(bot *botApi.BotAPI, postgresDb *config.PostgresDb) *deleteHandler {
	return &deleteHandler{
		command:  "delete",
		bot:      bot,
		database: postgresDb.Database,
	}
}
