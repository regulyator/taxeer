package config

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

type PostgresDb struct {
	Database *sql.DB
}

func InitConnection(host, user, password string) (*PostgresDb, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s?sslmode=disable", user, password, host))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresDb{Database: db}, nil
}

func RunMigration(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:./db/migration",
		"postgres", driver)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
}
