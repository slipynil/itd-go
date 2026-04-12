package dto

// AuthResponse представляет ответ от /api/v1/auth/refresh.
// Возвращает access token для аутентификации запросов.
type AuthResponse struct {
	// AccessToken - JWT токен для доступа к API
	AccessToken string `json:"accessToken"`
}

// UserIDResponse представляет ответ от /api/users/me с ID пользователя.
type UserIDResponse struct {
	// ID - уникальный идентификатор текущего пользователя
	ID string `json:"id"`
}
