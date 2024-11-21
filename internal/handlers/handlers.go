package handlers

import (
	"context"
	"fmt"
	"music/pkg/log"
	"net/http"
	"strconv"
	"strings"

	"music/internal/crud"
	"music/pkg/db"
	"music/pkg/server"
)

func GetLib(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("p")
    filter := r.URL.Query().Get("name")

    if pageStr == "" {
        pageStr = "1"
    }

    page, err := strconv.Atoi(pageStr)

    if err != nil {
        log.Logger.Debug(fmt.Sprintf("Неверный формат p=%s", pageStr))
        server.AnswerHandler(w, 400, "Неверный формат ввода страницы")
        return
    }

    list := []string{}
    out := make(chan string)
    answer := make(chan *[]string)
    go func() {
        for i := range out {
            
            if i == "Ошибка" {
                server.AnswerHandler(w, 500, "Ошибка запроса к базе данных")
                return
            }

            list = append(list, i)
        }

        answer<- &list
        close(answer)
    }()

    go db.DB.Query("list", out, crud.QueryList{
        P: page,
        Filter: filter,
    })

    server.AnswerHandler(w, 200, <-answer)
    log.Logger.Info("Запрос на получение данных библиотеки успешно выполнен")
}

func GetText(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)

    if err != nil {
        log.Logger.Debug(fmt.Sprintf("Неверный формат p=%s", idStr))
        server.AnswerHandler(w, 400, "Неверный формат ввода страницы")
        return
    }

    pageStr := r.URL.Query().Get("p")

    if pageStr == "" {
        pageStr = "1"
    }

    ctx := context.Background()
    rout := make(chan string)
    go db.GetKey(ctx, idStr, rout)

    page, err := strconv.Atoi(pageStr)

    if err != nil {
        log.Logger.Debug(fmt.Sprintf("Неверный формат p=%s", pageStr))
        server.AnswerHandler(w, 400, "Неверный формат ввода страницы")
        return
    }

    out := make(chan string)
    text := <- rout
    if text != "" {
        server.AnswerHandler(w, 200, strings.Split(text, "\n\n")[page-1])
        log.Logger.Debug(fmt.Sprintf("Текст песни с id=%d получен из кэша", id))
        return
    }

    answer := make(chan []string)
    go func() {
        var text []string
        for str := range out {
            if str == "Ошибка" {
                server.AnswerHandler(w, 400, "Песни с таким идентификатором не существует")
                log.Logger.Debug(fmt.Sprintf("Песни с id=%d не существует", id))
                return
            }
            
            text = append(text, str)
        }

        answer<- text
        close(answer)
        log.Logger.Debug(fmt.Sprintf("Текст песни с id=%d получен", id))
    }()

    go db.DB.Query("text", out, crud.Song{
        Id: id,
    })

    rout = make(chan string)
    go func() {
        res := <-rout

        if res != "success" {
            log.Logger.Error("Ошибка записи в кэш")
        }
    }()

    textPage := <-answer
    go db.NewKey(ctx, idStr, strings.Join(textPage, "\n\n"), rout)
    server.AnswerHandler(w, 200, textPage[page-1])
    log.Logger.Debug(fmt.Sprintf("Текст песни с id=%d получен", id))
}