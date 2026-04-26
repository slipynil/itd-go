package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError представляет ошибку, возвращённую ITD API.
// Содержит код ошибки, сообщение и HTTP статус код.
type APIError struct {
	// Code - код ошибки от API (например, "invalid_token", "not_found")
	Code string

	// Message - человекочитаемое описание ошибки
	Message string

	// StatusCode - HTTP статус код ответа
	StatusCode int
}

// Error возвращает строковое представление ошибки API.
// Реализует интерфейс error.
func (e *APIError) Error() string {
	return fmt.Sprintf("code: %s, message: %s, HTTP status: %d",
		e.Code, e.Message, e.StatusCode)
}

// CheckResponse проверяет HTTP статус код ответа на наличие ошибки.
// Возвращает APIError если статус >= 400, иначе nil.
// Поддерживает несколько форматов ошибок от ITD API.
func CheckResponse(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	// Пробуем формат 1: {"error": {"code": "...", "message": "..."}}
	var respErr ResponseError
	if err := json.Unmarshal(bodyBytes, &respErr); err == nil && respErr.Error.Code != "" {
		return &APIError{
			Code:       respErr.Error.Code,
			Message:    respErr.Error.Message,
			StatusCode: resp.StatusCode,
		}
	}

	// Пробуем формат 2: {"message": "...", "detail": "..."}
	var valErr ValidationError
	if err := json.Unmarshal(bodyBytes, &valErr); err == nil && valErr.Message != "" {
		return &APIError{
			Code:       "validation_error",
			Message:    valErr.Message + " " + valErr.Detail,
			StatusCode: resp.StatusCode,
		}
	}

	// Если не удалось распарсить, возвращаем сырое тело
	contentType := resp.Header.Get("Content-Type")
	var message string
	if !strings.Contains(contentType, "application/json") {
		message = fmt.Sprintf("server error (HTTP %d)", resp.StatusCode)
	} else {
		message = string(bodyBytes)
	}
	return &APIError{
		Code:       "unknown_error",
		Message:    message,
		StatusCode: resp.StatusCode,
	}
}
