package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"
)

// Извлекает домен из URL запроса
func extractDomain(request string) string {
	// Удалить порт, если он есть
	if idx := strings.Index(request, ":"); idx != -1 {
		request = request[:idx]
	}
	return request
}

// Проверяет, разрешен ли сайт, включая суб-домены
func isSiteAllowed(request string, allowedSites []string) bool {
	requestDomain := extractDomain(request)
	for _, site := range allowedSites {
		if requestDomain == site || strings.HasSuffix(requestDomain, "."+site) {
			return true // Возвращается true, если запрос соответствует одному из доменов или его субдоменам
		}
	}
	return false
}

// Проверяет, разрешен ли пользователь делать запросы
func isUserAllowed(username string, allowedUsers []string) bool {
	for _, user := range allowedUsers {
		if username == strings.ToLower(user) {
			return true
		}
	}
	return false
}

func getWindowsUsername() string {
	username := os.Getenv("USERNAME")
	if username == "" {
		return "Не удалось определить имя пользователя"
	}
	return strings.ToLower(username)
}

// Обрабатывает каждое подключение
func handleConnection(src net.Conn) {
	defer func() {
		if err := src.Close(); err != nil {
			logger("Ошибка при закрытии соединения: %v", err)
		}
	}()

	username := getWindowsUsername()
	configLock.RLock()
	allowedSitesArray := config.Allowed
	allowedUsers := config.Users
	configLock.RUnlock()

	buffer := make([]byte, 4096)
	n, err := src.Read(buffer)
	if err != nil {
		logger("Ошибка чтения из источника: %v", err)
		return
	}

	request := string(buffer[:n])
	if strings.HasPrefix(request, "CONNECT") {
		destination := request[len("CONNECT "):strings.Index(request, " HTTP/")]

		// Измененная логика: если пользователь не разрешен, позволяет доступ ко всем сайтам
		if isUserAllowed(username, allowedUsers) {
			// Если пользователь в списке разрешенных
			if isSiteAllowed(destination, allowedSitesArray) {
				// Сайт разрешен для данного пользователя
				logger("Запрос: %s\nПользователь: %s\nПользователь в списке: %v\nСайт разрешен: %v",
					destination,
					username,
					true,
					true)
			} else {
				// Сайт не разрешен для данного пользователя
				logger("Запрос: %s\nПользователь: %s\nПользователь в списке: %v\nСайт разрешен: %v",
					destination,
					username,
					true,
					false)
				return
			}
		} else {
			// Пользователь не в списке разрешенных, доступ ко всем сайтам
			logger("Запрос: %s\nПользователь: %s\nПользователь в списке: %v\nСайт разрешен: %v",
				destination,
				username,
				false,
				true)
		}

		dst, err := net.Dial("tcp", destination)
		if err != nil {
			logger("Не удалось соединиться с назначением: %v", err)
			return
		}
		defer func() {
			if err := dst.Close(); err != nil {
				logger("Ошибка при закрытии целевого соединения: %v", err)
			}
		}()

		// Отправляет клиенту подтверждение установки соединения
		if _, err := src.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n")); err != nil {
			logger("Ошибка при отправке подтверждения соединения: %v", err)
			return
		}

		// Начинает асинхронное копирование данных от клиента к серверу
		go func() {
			if _, err := io.Copy(dst, src); err != nil {
				logger("Ошибка при копировании данных к серверу: %v", err)
			}
		}()

		// Копирует данные от сервера к клиенту
		if _, err := io.Copy(src, dst); err != nil {
			logger("Ошибка при копировании данных к клиенту: %v", err)
		}
	}
}

func handleConnectionSocket(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			logger("Ошибка при закрытии соединения: %v", err)
		}
	}()
	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		logger("Ошибка чтения из соединения: %v", err)
		return
	}
	username = strings.TrimSpace(username)

	lock.Lock()
	currentUser = username
	activeUsers[username] = true
	lock.Unlock()

	logger("Зарегистрирован активный пользователь: %s", username)
}
