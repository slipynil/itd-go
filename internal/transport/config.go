package transport

import (
	"net/http"
	"time"

	"github.com/slipynil/itd-go/internal/auth"
)

// Config содержит конфигурацию для transport клиента.
type Config struct {
	// BaseURL - базовый URL для всех API запросов
	BaseURL string

	// HttpClient - настроенный HTTP клиент с middleware
	HttpClient *http.Client

	// AuthClient - провайдер аутентификации для добавления токенов к запросам
	AuthClient auth.Provider

	// MaxRetries - максимальное количество повторных попыток при ошибке 429
	MaxRetries int

	// RetryDelay - начальная задержка перед первой повторной попыткой
	RetryDelay time.Duration
}
