package transport

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	itderrors "github.com/slipynil/itd-go/errors"
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
	if err := itderrors.CheckResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// retryMiddleware автоматически повторяет запросы при ошибке 429 (rate limiting).
// Использует exponential backoff для задержки между попытками.
type retryMiddleware struct {
	base       http.RoundTripper
	maxRetries int
	baseDelay  time.Duration
}

// RoundTrip реализует интерфейс http.RoundTripper для retryMiddleware.
func (m *retryMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= m.maxRetries; attempt++ {
		// Выполняем запрос
		resp, err := m.base.RoundTrip(req)

		// Если нет ошибки или это не APIError, возвращаем результат
		if err == nil {
			return resp, nil
		}

		// Проверяем, является ли ошибка rate limiting (429)
		if !errors.Is(err, itderrors.ErrRateLimited) {
			return resp, err
		}

		lastErr = err

		// Если это последняя попытка, возвращаем ошибку
		if attempt == m.maxRetries {
			break
		}

		// Вычисляем задержку с exponential backoff
		delay := m.baseDelay * time.Duration(1<<uint(attempt))

		fmt.Printf("[itd-go] Rate limit exceeded (429), retry %d/%d after %v\n",
			attempt+1, m.maxRetries, delay)

		// Ждём с учётом контекста
		select {
		case <-time.After(delay):
			// Продолжаем следующую попытку
		case <-req.Context().Done():
			// Контекст отменён, прерываем retry
			return nil, req.Context().Err()
		}
	}

	return nil, lastErr
}
