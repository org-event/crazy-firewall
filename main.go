package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	configURL      = "https://org-event.github.io/crazy-firewall/config.json" // URL конфигурационного файла
	updateInterval = 5 * time.Minute                                          // Интервал обновления конфигурации
)

type Config struct {
	Users   []string `json:"users"`   // Список разрешенных пользователей
	Allowed []string `json:"allowed"` // Список разрешенных сайтов
}

var (
	config     Config            // Переменная для хранения конфигурации
	configLock = &sync.RWMutex{} // Мьютекс для синхронизированного доступа к конфигурации
)

// Загружает и обновляет конфигурацию каждые 5 минут
func fetchConfig() {
	for {
		resp, err := http.Get(configURL) // Выполнение HTTP GET запроса
		if err != nil {
			log.Printf("Не удалось получить конфигурацию: %v", err)
			time.Sleep(updateInterval)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)      // Чтение тела ответа
		if cerr := resp.Body.Close(); cerr != nil { // Немедленное закрытие тела ответа
			log.Printf("Не удалось закрыть тело ответа: %v", cerr)
		}
		if err != nil {
			log.Printf("Не удалось прочитать тело ответа: %v", err)
			time.Sleep(updateInterval)
			continue
		}

		var newConfig Config
		if err := json.Unmarshal(body, &newConfig); err != nil {
			log.Printf("Не удалось разобрать конфигурацию: %v", err)
			time.Sleep(updateInterval)
			continue
		}

		// Приведение имен пользователей к нижнему регистру
		for i, user := range newConfig.Users {
			newConfig.Users[i] = strings.ToLower(user)
		}

		configLock.Lock() // Блокирование доступа к конфигурации
		config = newConfig
		configLock.Unlock() // Разблокирование доступа

		log.Println("Конфигурация обновлена.")
		time.Sleep(updateInterval)
	}
}

// Проверяет, разрешен ли сайт, включая субдомены
func isBlockedSite(request string, allowedSites []string) bool {
	requestDomain := extractDomain(request)
	for _, site := range allowedSites {
		if requestDomain == site || strings.HasSuffix(requestDomain, "."+site) {
			return true // Возвращается true, если запрос соответствует одному из разрешенных доменов или его субдоменам
		}
	}
	return false
}

// Извлекает домен из URL запроса
func extractDomain(request string) string {
	// Обычно URL запроса приходит в формате `subdomain.domain.com:port`
	// Удалить порт, если он есть
	if idx := strings.Index(request, ":"); idx != -1 {
		request = request[:idx]
	}
	// Проверить наличие поддоменов и отрезать их
	if idx := strings.LastIndex(request, "."); idx != -1 {
		secondIdx := strings.LastIndex(request[:idx], ".")
		if secondIdx != -1 {
			request = request[secondIdx+1:]
		}
	}
	return request
}

// Проверяет, разрешен ли пользователь делать запросы
func isUserAllowed(username string, allowedUsers []string) bool {
	for _, user := range allowedUsers {
		if username == user {
			return true
		}
	}
	return false
}

// Обрабатывает каждое подключение
func handleConnection(src net.Conn) {
	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения: %v", err)
		}
	}()

	username := strings.ToLower(os.Getenv("USERNAME"))
	configLock.RLock()
	blockedSites := config.Allowed
	allowedUsers := config.Users
	configLock.RUnlock()

	buffer := make([]byte, 4096)
	n, err := src.Read(buffer)
	if err != nil {
		log.Printf("Ошибка чтения из источника: %v", err)
		return
	}

	request := string(buffer[:n])
	if strings.HasPrefix(request, "CONNECT") {
		destination := request[len("CONNECT "):strings.Index(request, " HTTP/")]
		if !isUserAllowed(username, allowedUsers) || !isBlockedSite(destination, blockedSites) {
			log.Printf("Запрос заблокирован: %s", destination)
			return
		}

		log.Printf("Запрос разрешен: %s", destination)
		dst, err := net.Dial("tcp", destination)
		if err != nil {
			log.Printf("Не удалось соединиться с назначением: %v", err)
			return
		}
		defer func() {
			if err := dst.Close(); err != nil {
				log.Printf("Ошибка при закрытии целевого соединения: %v", err)
			}
		}()

		// Отправляет клиенту подтверждение установки соединения
		if _, err := src.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n")); err != nil {
			log.Printf("Ошибка при отправке подтверждения соединения: %v", err)
			return
		}

		// Начинает асинхронное копирование данных от клиента к серверу
		go func() {
			if _, err := io.Copy(dst, src); err != nil {
				log.Printf("Ошибка при копировании данных к серверу: %v", err)
			}
		}()

		// Копирует данные от сервера к клиенту
		if _, err := io.Copy(src, dst); err != nil {
			log.Printf("Ошибка при копировании данных к клиенту: %v", err)
		}
	}
}

// Основная функция программы
func main() {
	go fetchConfig() // Запуск процесса обновления конфигурации

	listener, err := net.Listen("tcp", ":20000")
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}
	log.Println("Прослушивание порта :20000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Не удалось принять соединение: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
