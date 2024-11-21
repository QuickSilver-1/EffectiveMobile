package db

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	"music/internal/config"
	"music/pkg/log"

	_ "github.com/lib/pq"
)

var(
    DB = NewDB()
)

type ConnectDatabase struct {
    Quit        chan interface{}
    Connection  *sql.DB
    Command     map[string]func(*sql.DB, chan string, interface{})
}

// Query выполняет команду над базой данных
func (c *ConnectDatabase) Query(comm string, out chan string, data interface{}) {
    log.Logger.Debug(fmt.Sprintf("Выполнение команды: %s", comm))
    c.Command[comm](c.Connection, out, data)
}

// NewDB создаёт новое подключение к базе данных
func NewDB() ConnectDatabase {
    flag.Parse()

    log.Logger.Info("Подключение к базе данных")
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
                            config.AppConfig.PgHost, config.AppConfig.PgPort, config.AppConfig.PgUser, config.AppConfig.PgPass, config.AppConfig.PgName)
    conn, err := sql.Open("postgres", psqlInfo)

    if err != nil {
        log.Logger.Fatal(fmt.Sprintf("Проблема с подключением к базе данных: %v", err))
        panic("Проблема с подключением к базе данных")
    }

    conn.SetMaxOpenConns(10)
    conn.SetMaxIdleConns(5)
    conn.SetConnMaxLifetime(time.Second * 3)

    db := ConnectDatabase{
        Quit: make(chan interface{}),
        Connection: conn,
        Command: make(map[string]func(*sql.DB, chan string, interface{})),
    }

    log.Logger.Info("Успешное подключение к базе данных")

    go func() {
        <-db.Quit
        db.Connection.Close()
        log.Logger.Info("Подключение к базе данных закрыто")
    }()

    return db
}
