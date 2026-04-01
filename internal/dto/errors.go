package dto

// структура Error представляет ошибку сервера
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
