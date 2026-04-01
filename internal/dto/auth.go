package dto

// authResponse представляет ответ от /api/v1/auth/refresh
// возвращает access token
type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}

// UserMeResponse представляет ответ от /api/users/me
type UserIDResponse struct {
	ID string `json:"id"`
}
