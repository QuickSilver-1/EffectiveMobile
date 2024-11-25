package db

import (
	"database/sql"
	"fmt"
	"music/pkg/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// CreateSchema выполняет миграции базы данных для создания схемы
func CreateSchema(db *sql.DB) {
    log.Logger.Debug("Начало миграции")

    // Создаем экземпляр драйвера для PostgreSQL
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Logger.Fatal(fmt.Sprintf("Ошибка создания драйвера PostgreSQL: %v", err))
		panic("Ошибка создания драйвера PostgreSQL")
    }

    // Создаем мигратор с указанным источником миграций и базой данных
    m, err := migrate.NewWithDatabaseInstance("file://../../internal/migrations", "postgres", driver)
    if err != nil {
		log.Logger.Fatal(fmt.Sprintf("Ошибка создания мигратора: %v", err))
		panic("Ошибка создания мигратора")
    }

    // Применяем миграции к базе данных
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Logger.Fatal(fmt.Sprintf("Ошибка применения миграций: %v", err))
		panic("Ошибка применения миграций")
    }
    
    log.Logger.Debug("Миграции успешно применены!")
}
