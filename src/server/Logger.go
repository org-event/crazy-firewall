package main

import (
	"log"
	"os"
)

const (
	debugging  = true                            // Если true, сообщения выводятся в консоли
	logOutFile = true                            // Если true, вывод в файл
	logPath    = "C:\\crazyFirewall\\server.log" // Куда будем логировать
)

// Логирование сообщений, если будет включено
func logger(format string, v ...interface{}) {
	if debugging {
		log.Printf(format, v...)
	}
}

// Вызывается автоматически перед вызовом main()
func init() {
	if logOutFile {
		// Проверяем и создаем директорию, если необходимо
		if err := os.MkdirAll(logPath[:len(logPath)-len("crazyFirewall.log")], os.ModePerm); err != nil {
			log.Fatal("Не удалось создать директории для файла лога:", err)
		}

		// Создаем файл лога
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal("Не удалось открыть файл лога:", err)
		}
		log.SetOutput(logFile) // Переопределение вывода логов
	}
}
