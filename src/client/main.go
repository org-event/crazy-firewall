package main

/*
	Установка клиента для запуска при входе пользователя.
	Положить исполняемый файл в удобном месте.
	Создайте задачу в планировщике задач Windows для запуска этого файла при входе каждого пользователя.
*/

import (
	"net"
	"os/user"
	"strings"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		logger("Socket Client Failed to get current user:", err)
		return
	}

	// Получаем только имя пользователя, исключая домен
	username := currentUser.Username
	if idx := strings.LastIndex(username, "\\"); idx != -1 {
		username = username[idx+1:]
	}

	conn, err := net.Dial("tcp", "localhost:20001")
	if err != nil {
		logger("Socket Client Failed to connect to server:", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// Отправляем только имя пользователя
	_, err = conn.Write([]byte(username + "\n"))
	if err != nil {
		logger("Socket Client Failed to send username:", err)
		return
	}
}
