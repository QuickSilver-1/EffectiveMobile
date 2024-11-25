package crud

import (
	"database/sql"
	"fmt"
	"music/pkg/log"
	"strings"
)

// getListDB получает список песен из базы данных
func getListDB(database *sql.DB, out chan string, data interface{}) {
    defer close(out)

    dataType := data.(QueryList)
    log.Logger.Debug("Получение списка песен из БД")
    rows, err := database.Query(` SELECT song_id, song.name, author.name, release, link FROM song JOIN author ON song.author = author.author_id WHERE song.name ILIKE '%' || $1 || '%' OR author.name ILIKE '%' || $1 || '%' OFFSET $2 LIMIT 10 `, dataType.Filter, 10*(dataType.P-1))

    if err != nil {
        log.Logger.Error(fmt.Sprintf("Ошибка запроса к БД: %v", err))
        out <- "Error"
        return
    }

	// Валидация данных, чтобы обоаботать пустые поля
    for rows.Next() {
        var idI, songI, authorI, releaseI, linkI interface{}
        err = rows.Scan(&idI, &songI, &authorI, &releaseI, &linkI)

        id := validationInt(idI)
        song := validationStr(songI)
        author := validationStr(authorI)
        release := validationTime(releaseI).Format("2006-01-02")
        link := validationStr(linkI)

        if err != nil {
            out <- "Error"
            log.Logger.Error(fmt.Sprintf("Ошибка при чтении строки: %v", err))
            return
        }

        out <- fmt.Sprintf("%d Песня: %s исполнителя: %s. Дата выхода: %s. Ссылка: %s", id, song, author, release, link)
    }

    log.Logger.Debug("Вывод данных библиотеки прошёл успешно")
}

// getText получает текст песни из базы данных по её ID
func getText(database *sql.DB, out chan string, data interface{}) {
    defer close(out)

    dataType := data.(Song)
    var text string
    log.Logger.Debug(fmt.Sprintf("Получение текста песни с id=%d из БД", dataType.Id))
    database.QueryRow(` SELECT "text" FROM song WHERE "song_id" = $1 `, dataType.Id).Scan(&text)

    if text == "" {
        log.Logger.Debug(fmt.Sprintf("Запрос с неверный id=%d", dataType.Id))
        out <- "Error"
        return
    }

    for _, i := range strings.Split(text, "\n\n") {
        out <- i
    }

    log.Logger.Debug(fmt.Sprintf("Вывод текста песни с id %d прошёл успешно", dataType.Id))
}

// delSong удаляет песню из базы данных по её ID
func delSong(database *sql.DB, out chan string, data interface{}) {
    defer close(out)

    dataType := data.(Song)
    log.Logger.Debug("Удаление песни из БД")
    res, err := database.Exec(` DELETE FROM song WHERE "song_id" = $1 `, dataType.Id)

    if err != nil {
        out <- err.Error()
        return
    }
    
    r, _ := res.RowsAffected()

    if r == 0 {
        out <- "Not exist"
        return
    }

    out <- "success"
    log.Logger.Debug(fmt.Sprintf("Удаление песни с id %d прошло успешно", dataType.Id))
}

// createSongDB создает новую песню в базе данных
func createSongDB(database *sql.DB, out chan string, data interface{}) {
    defer close(out)

    dataType := data.(Song)
    var id int
    log.Logger.Debug("Создание песни")
    database.QueryRow(` SELECT "id" FROM author WHERE "name" = $1 `, dataType.Author).Scan(&id)
    if id == 0 {
		// Создаем запись об авторе, если не добавлен ранее
        database.QueryRow(` INSERT INTO author ("name") VALUES ($1) RETURNING "author_id" `, dataType.Author).Scan(&id)
    }

    var song_id string
    err := database.QueryRow(` INSERT INTO song ("name", "author", "text", "release", "link") VALUES ($1, $2, $3, $4, $5) RETURNING "song_id" `, dataType.Name, id, dataType.Text, dataType.Release, dataType.Link).Scan(&song_id)

    if err != nil {
        out <- err.Error()
        return
    }

    out <- "success"
    log.Logger.Debug(fmt.Sprintf("Создание песни с id %d - %s прошло успешно", id, dataType.Name))
}

// changeSongDB обновляет данные песни в базе данных
func changeSongDB(database *sql.DB, out chan string, data interface{}) {
    defer close(out)

    log.Logger.Debug("Обновление данных")
    dataType := data.(Song)
    var id int
    database.QueryRow(` SELECT "id" FROM author WHERE "name" = $1 `, dataType.Author).Scan(&id)
    
    if id == 0 && dataType.Author != "" {
		// Создаем запись об авторе, если не добавлен ранее
        database.QueryRow(` INSERT INTO author ("name") VALUES ($1) RETURNING "author_id" `, dataType.Author).Scan(&id)
    }

    query := `UPDATE song SET `
    params := []interface{}{}
    paramCounter := 1
    
	// Динамически формируем запрос, чтобы изменить только те поля, которые были затронуты
    if dataType.Name != "" {
        query += fmt.Sprintf("name = $%d, ", paramCounter)
        params = append(params, dataType.Name)
        paramCounter++
    }
    
    if id != 0 {
        query += fmt.Sprintf("author = $%d, ", paramCounter)
        params = append(params, id)
        paramCounter++ 
    }
        
    if dataType.Text != "" {
        query += fmt.Sprintf("text = $%d, ", paramCounter)
        params = append(params, dataType.Text)
        paramCounter++ 
    }
    
    if dataType.Release != "" {
        query += fmt.Sprintf("release = $%d, ", paramCounter)
        params = append(params, dataType.Release)
        paramCounter++
    }
    
    if dataType.Link != "" {
        query += fmt.Sprintf("link = $%d, ", paramCounter)
        params = append(params, dataType.Link)
        paramCounter++
    } 
    
    query = query[:len(query)-2]
    query += fmt.Sprintf(" WHERE song_id = $%d", paramCounter)
    params = append(params, dataType.Id) 
    res, err := database.Exec(query, params...) 
    
    if err != nil {
        out <- err.Error()
        return
    }
    
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        out <- err.Error()
        return
    } 
    
    if rowsAffected == 0 {
        out2 := make(chan string)
		// Если id не существует, то создаем новую запись с полученными данными
        go createSongDB(database, out2, data)
        res := <-out2
        if res != "success" {
            out <- res
            return
        }
    }

    out <- "success"
    log.Logger.Debug(fmt.Sprintf("Обновление песни с id %d - %s прошло успешно", dataType.Id, dataType.Name))
}
