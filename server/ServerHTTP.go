package main

import (
	"log"
	"net"
)

func startHTTPListener() {
	listener, err := net.Listen("tcp", portFirewall)
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}
	logger("Прослушивание порта: %v", portFirewall)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger("Не удалось принять соединение: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
