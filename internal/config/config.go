package config

import (
	"fmt"
	"music/pkg/log"
	"os"

	"github.com/joho/godotenv"
)

var (
    AppConfig = NewConfig()
    SecretKey = []byte("M8axEo25vLElQ8n85CvmFRmNrFWt0YQq")
)

type Config struct {
    HttpPort    string
    PgHost      string
    PgPort      string
    PgName      string
    PgUser      string
    PgPass      string
    RedHost     string
    RedPort     string
    RedPass     string
}

// NewConfig создает и возвращает новый конфигурационный объект
func NewConfig() *Config {
    // Загрузка переменных окружения из .env файла
    err := godotenv.Load("../../config.env")
    
    if err != nil {
        fmt.Println(err)
        log.Logger.Fatal("Ошибка загрузки .env файла")
    }
    
    config := &Config{
        HttpPort: os.Getenv("HTTP_PORT"),   // Порт для запуска приложения
        PgHost: os.Getenv("DB_HOST"),       // Хост для базы данных PostgreSQL
        PgPort: os.Getenv("DB_PORT"),       // Порт для базы данных PostgreSQL
        PgName: os.Getenv("DB_NAME"),       // Имя базы данных PostgreSQL
        PgUser: os.Getenv("DB_USER"),       // Имя пользователя для базы данных PostgreSQL
        PgPass: os.Getenv("DB_PASSWORD"),   // Пароль для базы данных PostgreSQL
        RedHost: os.Getenv("RED_HOST"),     // Хост для базы данных Redis
        RedPort: os.Getenv("RED_PORT"),     // Порт для базы данных Redis
        RedPass: os.Getenv("RED_PASSWORD"), // Пароль для базы данных Redis
    }

    return config
}