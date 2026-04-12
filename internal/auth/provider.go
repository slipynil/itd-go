package auth

import (
	"context"
	"net/http"
)

// Provider определяет интерфейс для провайдера аутентификации.
// Предоставляет методы для получения токенов и управления сессией.
type Provider interface {
	// GetCookieJar возвращает хранилище cookies для HTTP клиента
	GetCookieJar() http.CookieJar

	// GetAccessToken получает актуальный access token, обновляя его при необходимости
	GetAccessToken(ctx context.Context) (string, error)

	// GetUserID возвращает ID текущего аутентифицированного пользователя
	GetUserID() string

	// IsAuthenticated проверяет, аутентифицирован ли пользователь
	IsAuthenticated() bool
}
