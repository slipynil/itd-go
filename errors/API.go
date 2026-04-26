package errors

import "errors"

// API errors - ошибки от сервера ITD

// ErrUnauthorized возвращается при невалидном/истёкшем токене (401).
var ErrUnauthorized = errors.New("unauthorized: invalid or expired token")

// ErrNotFound возвращается когда ресурс не найден (404).
var ErrNotFound = errors.New("resource not found")

// ErrForbidden возвращается при отсутствии прав доступа (403).
var ErrForbidden = errors.New("forbidden: insufficient permissions")

// ErrRateLimited возвращается при превышении лимита запросов (429).
var ErrRateLimited = errors.New("rate limit exceeded")

// ErrServerError возвращается при внутренней ошибке сервера (500+).
var ErrServerError = errors.New("internal server error")

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
