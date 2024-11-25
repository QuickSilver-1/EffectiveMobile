package crud

import (
	"database/sql"

	"music/pkg/db"
	"music/pkg/log"
)

// CollectHandlers инициализирует обработчики CRUD операций
func CollectHandlers(conn *db.ConnectDatabase) {
    log.Logger.Info("Сборка CRUD")

    // Инициализация команд в структуре ConnectDatabase
    conn.Command = map[string]func(*sql.DB, chan string, interface{}){
        "list":     getListDB,
        "text":     getText,
        "delete":   delSong,
        "new":      createSongDB,
        "change":   changeSongDB,
    }
}
