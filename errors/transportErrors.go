package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (e *APIError) Error() string {
	return fmt.Sprintf("код: %s, сообщение: %s, HTTP статус: %d",
		e.Code, e.Message, e.StatusCode)
}

// проверяет http статус кода на наличие ошибки и возвращает ее
func CheckResponse(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("не удалось прочитать тело ответа: %w", err)
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
	return &APIError{
		Code:       "unknown_error",
		Message:    string(bodyBytes),
		StatusCode: resp.StatusCode,
	}
}
