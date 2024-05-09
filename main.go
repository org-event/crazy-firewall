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

		body, err := ioutil.ReadAll(resp.Body) // Чтение тела ответа
		resp.Body.Close()
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

// Проверяет, разрешен ли сайт
func isBlockedSite(request string, blockedSites []string) bool {
	for _, site := range blockedSites {
		if strings.Contains(request, site) {
			return true
		}
	}
	return false
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
	defer src.Close()

	username := strings.ToLower(os.Getenv("USERNAME")) // Получение имени пользователя в нижнем регистру
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
		defer dst.Close()

		src.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		go io.Copy(dst, src)
		io.Copy(src, dst)
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
