package crud

import (
	"database/sql"
	"fmt"
	"music/pkg/log"
	"strings"
)

func getListDB(database *sql.DB, out chan string, data interface{}) {
	defer close(out)

	dataType := data.(QueryList)
	log.Logger.Debug("Получение списка песен из БД")
	rows, err := database.Query(` SELECT song_id, song.name, author.name, release, link FROM song JOIN author ON song.author = author.author_id WHERE song.name ILIKE '%' || $1 || '%' OR author.name ILIKE '%' || $1 || '%' OFFSET $2 LIMIT 10 `, dataType.Filter, 10*(dataType.P-1))

	if err != nil {
		log.Logger.Error(fmt.Sprintf("Ошибка запроса к БД: %v", err))
		out<- "Ошибка"
		return
	}

	for rows.Next() {
		var id, song, author, release, link string
		err = rows.Scan(&id, &song, &author, &release, &link)

		if err != nil {
			out<- "Ошибка"
            log.Logger.Error(fmt.Sprintf("Ошибка при чтении строки: %v", err))
            return
        }

		out<- fmt.Sprintf("%s Песня: %s исполнителя: %s. Дата выхода: %s. Ссылка: %s", id, song, author, release, link)
	}

	log.Logger.Debug("Вывод данных библиотеки прошёл успешно")
}

func getText(database *sql.DB, out chan string, data interface{}) {
	defer close(out)

	dataType := data.(Song)
	var text string
	log.Logger.Debug(fmt.Sprintf("Получение текста песни с id=%d из БД", dataType.Id))
	database.QueryRow(` SELECT "text" FROM song WHERE "song_id" = $1 `, dataType.Id).Scan(&text)

	if text == "" {
		log.Logger.Debug(fmt.Sprintf("Запрос с неверный id=%d", dataType.Id))
		out<- "Ошибка"
		return
	}

	for _, i := range strings.Split(text, "\n\n") {
		out<- i
	}

	log.Logger.Debug(fmt.Sprintf("Вывод текста песни с id %d прошёл успешно", dataType.Id))
}