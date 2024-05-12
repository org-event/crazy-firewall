package main

import (
	"fmt"
	"net"
	"os/user"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user:", err)
		return
	}

	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	_, err = conn.Write([]byte(currentUser.Username + "\n"))
	if err != nil {
		fmt.Println("Failed to send username:", err)
		return
	}
}

/*
	Установка клиента для запуска при входе пользователя.
	Положить исполняемый файл в удобном месте.
	Создайте задачу в планировщике задач Windows для запуска этого файла при входе каждого пользователя.
*/
