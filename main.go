package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// Функция проверяет, содержится ли адрес в списке сайтов
func isBlockedSite(request string, blockedSites []string) bool {
	for _, site := range blockedSites {
		if strings.Contains(request, site) {
			return true
		}
	}
	return false
}

// Обработчик каждого подключения
func handleConnection(src net.Conn) {
	defer src.Close()

	username := os.Getenv("USERNAME")
	blockedSites := []string{"ya.ru", "vk.ru", "developer.mozilla.org", "motobavaria.com"}

	buffer := make([]byte, 4096)
	n, err := src.Read(buffer)
	if err != nil {
		log.Printf("Error reading from source: %v", err)
		return
	}

	request := string(buffer[:n])
	if strings.HasPrefix(request, "CONNECT") {
		destination := request[len("CONNECT "):strings.Index(request, " HTTP/")]
		if username == "Ivan" && !isBlockedSite(destination, blockedSites) {
			log.Printf("Запрос заблокирован: %s", destination)
			return
		}

		log.Printf("Запрос разрешен: %s", destination)
		dst, err := net.Dial("tcp", destination)
		if err != nil {
			log.Printf("Unable to connect to destination: %v", err)
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
	listener, err := net.Listen("tcp", ":20000")
	if err != nil {
		log.Fatalf("Unable to listen on port: %v", err)
	}
	log.Println("Listening on :20000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
