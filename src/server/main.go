package main

func main() {
	go fetchConfig() // Запуск процесса обновления конфигурации
	go startHTTPListener()
	go startUserListener()

	// Ожидаем завершения горутин, можно использовать waitgroup или подобное
	select {} // Этот select блокирует main функцию, предотвращая её завершение
}
