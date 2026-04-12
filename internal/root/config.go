package root

import (
	"time"
)

// Config содержит параметры конфигурации для корневого клиента SDK.
type Config struct {
	// RefreshToken - токен для обновления сессии
	RefreshToken string

	// Url - базовый URL API
	Url string

	// Domain - доменное имя API
	Domain string

	// UserAgent - User-Agent заголовок для HTTP запросов
	UserAgent string

	// Timeout - таймаут для HTTP запросов
	Timeout time.Duration
}
