package db

import (
	"database/sql"
	"music/pkg/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func CreateSchema(db *sql.DB) {
	log.Logger.Debug("Начало миграции")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	
	if err != nil {
		log.Logger.Fatal(err.Error()) }
		
	m, err := migrate.NewWithDatabaseInstance("file://../../internal/migrations", "postgres", driver)
	
	if err != nil {
		log.Logger.Fatal(err.Error()) }
		
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Logger.Fatal(err.Error())
	}
	
	log.Logger.Debug("Миграции успешно применены!")
}
