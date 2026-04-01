package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/slipynil/itd-go/internal/dto"
)

type APIError struct {
	Code       string
	Message    string
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
	data := new(dto.Error)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}
	return &APIError{
		Code:       data.Code,
		Message:    data.Message,
		StatusCode: resp.StatusCode,
	}
}
