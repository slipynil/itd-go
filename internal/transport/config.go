package transport

import (
	"net/http"

	"github.com/slipynil/itd-go/internal/auth"
)

// Config содержит конфигурацию transport клиента
type Config struct {
	BaseURL    string        // базовый URL
	HttpClient *http.Client  // авторизованный и настроенный [http.Client]
	AuthClient auth.Provider // структура провайдера аутентификации
}
