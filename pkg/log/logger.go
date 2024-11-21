package log

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var (
    Logger = loggerBuild()
)

// loggerBuild инициализирует и возвращает экземпляр логгера Zap
func loggerBuild() *zap.Logger {
    config := zap.NewDevelopmentConfig()
    config.OutputPaths = []string{"../../log/log.log", "stdout"}

    logger, err := config.Build()

    if err != nil {
        file, _ := os.OpenFile("../../log/log.log", os.O_APPEND, 0666)
        log.SetOutput(file)
        log.Fatal("Failed to open log file: ", err)
        panic("failed to configure logger")
    }

    return logger
}
