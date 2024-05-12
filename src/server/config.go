package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	portFirewall   = ":20000"
	portSocket     = ":20001"
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

	activeUsers = make(map[string]bool)
	lock        sync.Mutex
	currentUser string
)

// Загружает и обновляет конфигурацию каждые 5 минут
func fetchConfig() {
	for {
		resp, err := http.Get(configURL) // Выполнение HTTP GET запроса
		if err != nil {
			logger("Не удалось получить конфигурацию: %v", err)
			time.Sleep(updateInterval)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)      // Чтение тела ответа
		if cerr := resp.Body.Close(); cerr != nil { // Немедленное закрытие тела ответа
			logger("Не удалось закрыть тело ответа: %v", cerr)
		}
		if err != nil {
			logger("Не удалось прочитать тело ответа: %v", err)
			time.Sleep(updateInterval)
			continue
		}

		var newConfig Config
		if err := json.Unmarshal(body, &newConfig); err != nil {
			logger("Не удалось разобрать конфигурацию: %v", err)
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

		logger("Конфигурация обновлена.")
		time.Sleep(updateInterval)
	}
}
