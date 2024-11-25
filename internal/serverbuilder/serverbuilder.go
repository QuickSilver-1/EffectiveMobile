package serverbuilder

import (
	"net/http"

	"music/internal/handlers"
	"music/pkg/log"
	"music/pkg/server"

	"github.com/gorilla/mux"
)

// MakeServer создает и настраивает HTTP сервер с маршрутизацией
func MakeServer(port string, readWait, writeWait int) *http.Server {
    log.Logger.Info("Формирование рутера")
    mux := mux.NewRouter()

    // Устанавливаем middleware для логирования запросов и защиты от brute force
    mux.Use(server.Middleware)
    mux.Use(server.LimitMiddleware)

    mux.HandleFunc("/list", handlers.GetLib).Methods("GET")
    mux.HandleFunc("/text", handlers.GetText).Methods("GET")
    mux.HandleFunc("/del", handlers.DelSong).Methods("DELETE")
    mux.HandleFunc("/new", handlers.CreateSong).Methods("POST")
    mux.HandleFunc("/change", handlers.ChangeSong).Methods("PUT")
    return server.NewServer(port, mux, readWait, writeWait)
}
