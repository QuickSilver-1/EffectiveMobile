package main

import (
	"fmt"

	"music/internal/config"
	dh "music/internal/crud"
	sb "music/internal/serverbuilder"
	"music/pkg/db"
	"music/pkg/log"
)

func main() {
	// Запустить миграцию
	db.CreateSchema(db.DB.Connection)
	
    // Сборка обработчиков базы данных
    dh.CollectHandlers(&db.DB)

    // Cоздания сервера
    appServer := sb.MakeServer(":" + config.AppConfig.HttpPort, 10, 10)
    
    log.Logger.Info(fmt.Sprintf("Сервер запущен на порту %s", config.AppConfig.HttpPort))
    
    // Запуск сервера
    err := appServer.ListenAndServe()
    
    if err != nil {
        log.Logger.Error(fmt.Sprintf("Ошибка при запуске сервера: %v", err))
    }
}
