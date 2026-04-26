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

	// Err - обёрнутая sentinel ошибка для проверки через errors.Is()
	Err error
}

// Error возвращает строковое представление ошибки API.
// Реализует интерфейс error.
func (e *APIError) Error() string {
	return fmt.Sprintf("code: %s, message: %s, HTTP status: %d",
		e.Code, e.Message, e.StatusCode)
}

// Unwrap возвращает обёрнутую sentinel ошибку.
// Позволяет использовать errors.Is() для проверки типа ошибки.
func (e *APIError) Unwrap() error {
	return e.Err
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

	var code, message string

	// Пробуем формат 1: {"error": {"code": "...", "message": "..."}}
	var respErr ResponseError
	if err := json.Unmarshal(bodyBytes, &respErr); err == nil && respErr.Error.Code != "" {
		code = respErr.Error.Code
		message = respErr.Error.Message
	} else {
		// Пробуем формат 2: {"message": "...", "detail": "..."}
		var valErr ValidationError
		if err := json.Unmarshal(bodyBytes, &valErr); err == nil && valErr.Message != "" {
			code = "validation_error"
			message = valErr.Message + " " + valErr.Detail
		} else {
			// Если не удалось распарсить, возвращаем сырое тело
			contentType := resp.Header.Get("Content-Type")
			code = "unknown_error"
			if !strings.Contains(contentType, "application/json") {
				message = fmt.Sprintf("server error (HTTP %d)", resp.StatusCode)
			} else {
				message = string(bodyBytes)
			}
		}
	}

	apiErr := &APIError{
		Code:       code,
		Message:    message,
		StatusCode: resp.StatusCode,
	}

	// Маппинг HTTP статусов на sentinel errors
	switch {
	case resp.StatusCode == 401:
		apiErr.Err = ErrUnauthorized
	case resp.StatusCode == 403:
		apiErr.Err = ErrForbidden
	case resp.StatusCode == 404:
		apiErr.Err = ErrNotFound
	case resp.StatusCode == 429:
		apiErr.Err = ErrRateLimited
	case resp.StatusCode >= 500:
		apiErr.Err = ErrServerError
	}

	return apiErr
}
