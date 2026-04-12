package transport

import (
	"net/http"

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
}
