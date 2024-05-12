package main

import (
	"bufio"
	"io"
	"net"
	"strings"
)

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

// Проверяет, разрешен ли сайт, включая суб-домены
func isSiteDisallow(request string, allowedSites []string) bool {
	requestDomain := extractDomain(request)
	for _, site := range allowedSites {
		if requestDomain == site || strings.HasSuffix(requestDomain, "."+site) {
			return false // Возвращается false, если запрос соответствует одному из доменов или его субдоменам
		}
	}
	return true
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
			logger("Ошибка при закрытии соединения: %v", err)
		}
	}()

	username := currentUser
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

		if isUserAllowed(username, allowedUsers) && isSiteDisallow(destination, allowedSitesArray) {
			logger(
				"Запрос: %s",
				destination,
				"Пользователь ",
				username,
				"Пользователь в списке",
				isUserAllowed(username, allowedUsers),
				"сайт запрещен",
				isSiteDisallow(destination, allowedSitesArray))
			return
		}

		logger("Запрос разрешен: %s", destination)
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
		logger("Failed to read from connection: %v", err)
		return
	}
	username = strings.TrimSpace(username)

	lock.Lock()
	currentUser = username
	activeUsers[username] = true
	lock.Unlock()

	logger("Registered active user: %s", username)
}
