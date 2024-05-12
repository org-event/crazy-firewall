package main

import (
	"net"
	"os"
)

func startUserListener() {
	listenerSocket, err := net.Listen("tcp", portSocket)
	if err != nil {
		logger("Error starting server: %v", err)
		os.Exit(1)
	}

	defer func() {
		if err := listenerSocket.Close(); err != nil {
			logger("Ошибка при закрытии сокет соединения: %v", err)
		}
	}()

	logger("Server is listening on %v", portSocket)
	for {
		conn, err := listenerSocket.Accept()
		if err != nil {
			logger("Error accepting connection: %v", err)
			continue
		}
		go handleConnectionSocket(conn)
	}
}
