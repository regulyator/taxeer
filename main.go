package main

import (
	"log"
	"os"
	taxeerBot "taxeer/bot"
	"taxeer/util/config"
)

func main() {
	dbConnection, err := config.InitConnection(os.Getenv("DB_HOST_KEY"), os.Getenv("DB_USER_KEY"), os.Getenv("DB_PASSWORD_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	config.RunMigration(dbConnection.Database)
	taxeerBot.StartTaxeerBotListener(dbConnection)
}
