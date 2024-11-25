package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"music/pkg/log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"music/internal/crud"
	"music/pkg/db"
	"music/pkg/server"
)

// GetLib получает список песен из библиотеки
func GetLib(w http.ResponseWriter, r *http.Request) {
    pageStr := r.URL.Query().Get("p")
    filter := r.URL.Query().Get("name")

    if pageStr == "" {
        pageStr = "1"
    }

    page, err := strconv.Atoi(pageStr)

    if err != nil || page < 1 {
        log.Logger.Debug(fmt.Sprintf("Неверный формат p=%s", pageStr))
        server.AnswerHandler(w, 400, "Неверный формат ввода страницы")
        return
    }

    list := []string{}
    out := make(chan string)
    answer := make(chan *[]string)
    go func() {
        for i := range out {
            if i == "Error" {
                server.AnswerHandler(w, 500, "Ошибка запроса к базе данных")
                log.Logger.Error(fmt.Sprintf("Ошибка запроса к БД. Ошибка %s", i))
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

// GetText получает текст песни по её ID
func GetText(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)

    if err != nil {
        log.Logger.Debug(fmt.Sprintf("Неверный формат id=%s", idStr))
        server.AnswerHandler(w, 400, "Неверный формат ввода id")
        return
    }

    pageStr := r.URL.Query().Get("p")

    if pageStr == "" {
        pageStr = "1"
    }

    ctx := context.Background()
    rout := make(chan string)
    // Получаем данные из кэша
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
            if str == "Error" {
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
            log.Logger.Error(fmt.Sprintf("Ошибка записи в кэш. Ошибка: %s", res))
        }
    }()

    textPage := <-answer
    // Создаём запись в кэше
    go db.NewKey(ctx, idStr, strings.Join(textPage, "\n\n"), rout)
    server.AnswerHandler(w, 200, textPage[page-1])
    log.Logger.Debug(fmt.Sprintf("Текст песни с id=%d получен", id))
}

// DelSong удаляет песню по её ID
func DelSong(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")

    if idStr == "" {
        log.Logger.Debug("id не указан")
        server.AnswerHandler(w, 400, "id не указан")
        return
    }

    id, err := strconv.Atoi(idStr)

    if err != nil {
        log.Logger.Debug(fmt.Sprintf("Неверный формат id=%s", idStr))
        server.AnswerHandler(w, 400, "Неверный формат ввода id")
        return
    }

    wg := &sync.WaitGroup{}
    out := make(chan string)
    wg.Add(1)
    go func() {
        defer wg.Done()

        res := <-out
        if res == "Not exist" {
            server.AnswerHandler(w, 400, "Песни с таким id не существует")
            return
        }

        if res != "success" {
            server.AnswerHandler(w, 500, "Ошибка удаления песни")
            log.Logger.Error(fmt.Sprintf("Ошибка удаления песни с id %d. Ошибка: %s", id, res))
            return
        }
        
        server.AnswerHandler(w, 200, "Песня успешно удалена")
        log.Logger.Debug(fmt.Sprintf("Удаление песни с id %d прошло успешно", id))
    }()
    go db.DB.Query("delete", out, crud.Song{
        Id: id,
    })

    ctx := context.Background()

    rout := make(chan string)
    wg.Add(1)
    go func() {
        defer wg.Done()

        res := <-rout
        if res != "success" {
            log.Logger.Error(fmt.Sprintf("Ошибка удаления песни с id %d из Redis. Ошибка: %s", id, res))
            return
        }

        log.Logger.Debug(fmt.Sprintf("Удаление песни с id %d из кэша прошло успешно", id))
    }()
    // Удаляем запись из кэша
    go db.DelKey(ctx, idStr, rout)

    wg.Wait()
}

// CreateSong создает новую песню
func CreateSong(w http.ResponseWriter, r *http.Request) {
    var song crud.Song
    json.NewDecoder(r.Body).Decode(&song)
    date := strings.Split(song.Release, ".")
    song.Release = date[2] + "-" + date[1] + "-" + date[0]

    out := make(chan string)        

    go db.DB.Query("new", out, song)

    var resp *http.Response
    var err error
    // Формируем запрос на тот же адерс, где находится этот сервер
    if r.TLS == nil {
        resp, err = http.Get(fmt.Sprintf("http://%s/info?group=%s&song=%s", r.Host, song.Author, song.Name))
    } else {
        resp, err = http.Get(fmt.Sprintf("https://%s/info?group=%s&song=%s", r.Host, song.Author, song.Name))
    }

    if err != nil {
        log.Logger.Error(fmt.Sprintf("Ошибка запроса: %v", err))
    }

    res := <-out
    if res != "success" {
        server.AnswerHandler(w, 500, "Ошибка загрузки песни")
        log.Logger.Error(fmt.Sprintf("Ошибка создания песни с именем %s. Ошибка: %s", song.Name, res))
        return
    }

    server.AnswerHandler(w, 200, resp.Body)
    log.Logger.Debug("Создание новой записи прошло успешно")
}

// ChangeSong изменяет данные существующей песни
func ChangeSong(w http.ResponseWriter, r *http.Request) {
    // Валидируем данные
    idStr := r.URL.Query().Get("id")
    var song crud.Song
    json.NewDecoder(r.Body).Decode(&song)
    id, err := strconv.Atoi(idStr)
    song.Id = id
    date := strings.Split(song.Release, ".")
    song.Release = date[2] + "-" + date[1] + "-" + date[0]

    if err != nil {
        server.AnswerHandler(w, 400, "Неверный id")
        return
    }

    out := make(chan string)
    go func() {
        res := <-out

        if res != "success" {
            server.AnswerHandler(w, 500, "Ошибка изменения данных")
            log.Logger.Error(fmt.Sprintf("Ошибка изменения данных id: %d. Ошибка: %s", song.Id, res))
            return
        }

        server.AnswerHandler(w, 200, "Данные обновлены")
    }()

    db.DB.Query("change", out, song)
    log.Logger.Debug("Обновление записи прошло успешно")
}