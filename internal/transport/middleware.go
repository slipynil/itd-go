package transport

import (
	"net/http"

	"github.com/slipynil/itd-go/errors"
	"github.com/slipynil/itd-go/internal/auth"
)

// authMiddleware добавляет заголовок Authorization с Bearer токеном к каждому запросу.
type authMiddleware struct {
	base     http.RoundTripper // базовый транспорт для выполнения запроса
	provider auth.Provider     // провайдер аутентификации для получения токена
}

// RoundTrip реализует интерфейс http.RoundTripper для authMiddleware.
func (m *authMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	// Получаем токен от провайдера
	token, err := m.provider.GetAccessToken(req.Context())
	if err != nil {
		return nil, err
	}

	// Клонируем запрос и добавляем заголовок
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+token)

	return m.base.RoundTrip(req)
}

// statusCheckMiddleware проверяет HTTP статус код ответа и возвращает ошибку при 4xx/5xx.
type statusCheckMiddleware struct {
	base http.RoundTripper
}

// RoundTrip реализует интерфейс http.RoundTripper для statusCheckMiddleware.
func (m *statusCheckMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := m.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Проверяем статус ответа
	if err := errors.CheckResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
