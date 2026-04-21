package errors

// Error представляет структуру ошибки от сервера ITD API.
// Содержит код ошибки и человекочитаемое сообщение.
type Error struct {
	// Code - код ошибки (например, "invalid_token", "not_found")
	Code string `json:"code"`

	// Message - описание ошибки
	Message string `json:"message"`
}

// ResponseError представляет формат ошибки API: {"error": {"code": "...", "message": "..."}}.
// Используется для парсинга стандартных ошибок от ITD API.
type ResponseError struct {
	// Error - вложенная структура с деталями ошибки
	Error Error `json:"error"`
}

// ValidationError представляет формат ошибки валидации: {"message": "...", "detail": "..."}.
// Используется для парсинга ошибок валидации входных данных.
type ValidationError struct {
	// Message - основное сообщение об ошибке
	Message string `json:"message"`

	// Detail - дополнительные детали ошибки
	Detail string `json:"detail"`
}
