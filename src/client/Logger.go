package main

import (
	"log"
	"os"
	"path/filepath"
)

const (
	debugging  = true                            // Если true, сообщения выводятся в консоли
	logOutFile = true                            // Если true, вывод в файл
	logPath    = "C:\\crazyFirewall\\client.log" // Куда будет вывод
)

// Логирование сообщений, если будет включено
func logger(format string, v ...interface{}) {
	if debugging {
		log.Printf(format, v...)
	}
}

// Вызывается автоматически перед вызовом main()
// Вызывается автоматически перед вызовом main()
func init() {
	if logOutFile {
		// Проверяем и создаем директорию, если необходимо
		dirPath := filepath.Dir(logPath) // Получаем путь к директории из полного пути файла
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Fatal("Не удалось создать директории для файла лога:", err)
		}

		// Проверяем, существует ли файл
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			// Файла не существует, создаем его и добавляем BOM для UTF-8
			logFile, err := os.Create(logPath)
			if err != nil {
				log.Fatal("Не удалось создать файл лога:", err)
			}
			defer func(logFile *os.File) {
				err := logFile.Close()
				if err != nil {

				}
			}(logFile)
			if _, err = logFile.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil { // BOM для UTF-8
				log.Fatal("Не удалось записать BOM в файл лога:", err)
			}
		}

		// Открываем файл для записи (файл уже должен существовать)
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal("Не удалось открыть файл лога для дописывания:", err)
		}
		log.SetOutput(logFile) // Назначаем файл выводом для пакета log
	}
}
