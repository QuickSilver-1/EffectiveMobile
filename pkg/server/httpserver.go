package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"music/pkg/log"

	"github.com/gorilla/mux"
)

type Answer struct {
    StatusCode int
    Value      interface{}
}

// AnswerHandler обрабатывает и отправляет ответ клиенту
func AnswerHandler(w http.ResponseWriter, code int, value interface{}) {
    w.Header().Set("Content-Type", "application/json")

    answer := Answer{
        StatusCode: code,
        Value: value,
    }

    err := json.NewEncoder(w).Encode(answer)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Logger.Error(fmt.Sprintf("Ошибка кодирования ответа: %v", err))
        return
    }

    log.Logger.Debug(fmt.Sprintf("Ответ успешно обработан: %v", answer))
}

// NewServer создает новый HTTP сервер с заданными параметрами
func NewServer(port string, mux *mux.Router, r, w int) *http.Server {
    log.Logger.Info(fmt.Sprintf("Создание нового сервера на порту %s", port))

    server := &http.Server{
        Addr:         port,
        Handler:      mux,
        ReadTimeout:  time.Duration(r) * time.Second,
        WriteTimeout: time.Duration(w) * time.Second,
    }

    log.Logger.Info("Новый сервер успешно создан")
    return server
}

// StartServer запускает HTTP сервер и обрабатывает возможные ошибки
func StartServer(server *http.Server) error {
    log.Logger.Info("Запуск сервера")

    err := server.ListenAndServe()
    if err != nil {
        log.Logger.Error(fmt.Sprintf("Ошибка при запуске сервера: %v", err))
    }
    return err
}