package transport

import (
	"fmt"
	"net/http"

	"github.com/slipynil/itd-go/internal/auth"
	"github.com/slipynil/itd-go/internal/errors"
)

// authMiddleware добавляет заголовок Authorization к каждому запросу
type authMiddleware struct {
	base     http.RoundTripper // структура с реализацией интерфейса для HTTP клиента
	provider auth.Provider     // структура с интерфейсом провайдера аутентификации
}

// реализация метода RoundTrip для authMiddleware
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

// statusCheckMiddleware проверяет статус ответа и возвращает ошибку при необходимости
type statusCheckMiddleware struct {
	base http.RoundTripper
}

func (m *statusCheckMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := m.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	fmt.Println("Протокол: ", resp.Proto)

	// Проверяем статус ответа
	if err := errors.CheckResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}

	return resp, nil
}
