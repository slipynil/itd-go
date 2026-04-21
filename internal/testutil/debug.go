package testutil

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// RawAnswerAndStruct - выводит сырые данные в json string формат.
// Функция предназначена только для тестирования http ответов.
// Возвращает сырые данные и парсит данные по указателю в структуру.
func RawAnswerAndStruct(r io.ReadCloser, dto any, opts ...json.Options) (string, error) {

	if r == nil {
		return "", fmt.Errorf("reader is nil")
	}

	if reflect.ValueOf(dto).Kind() != reflect.Pointer {
		return "", fmt.Errorf("dto must be a pointer, got %T %v", dto, dto)
	}

	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	if err := json.UnmarshalRead(tee, dto, opts...); err != nil {
		return "", err
	}

	// Выводим сырой ответ от сервера
	rawResponse := buf.String()
	fmt.Printf("RAW SERVER RESPONSE: %s\n\n", rawResponse)
	fmt.Printf("dto after unmarshal: %+v\n", dto)

	prettyBytes, err := json.Marshal(dto, jsontext.WithIndent("   "))
	if err != nil {
		return buf.String(), nil
	}

	return string(prettyBytes), nil
}

// RawAnswer - выводит сырые данные в json string формат.
func RawAnswer(r io.ReadCloser) error {

	if r == nil {
		return fmt.Errorf("reader is nil")
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}
	rawResponse := buf.String()
	fmt.Printf("RAW SERVER RESPONSE: %s\n\n", rawResponse)
	return nil
}
