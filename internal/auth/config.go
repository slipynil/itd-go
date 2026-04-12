package auth

import "net/http"

// Config содержит параметры конфигурации для клиента аутентификации.
type Config struct {
	// Url - базовый URL для запросов аутентификации
	Url string

	// HttpClient - HTTP клиент для выполнения запросов
	HttpClient *http.Client
}
