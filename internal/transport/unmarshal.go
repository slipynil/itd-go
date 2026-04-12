package transport

import (
	"strings"
	"time"

	"github.com/go-json-experiment/json"
)

var timeUnmarshaler = json.UnmarshalFunc(func(b []byte, t *time.Time) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// Пробуем ISO8601 формат (RFC3339)
	parsed, err := time.Parse(time.RFC3339, s)
	if err == nil {
		*t = parsed
		return nil
	}

	// Пробуем формат с микросекундами
	layout := "2006-01-02 15:04:05.999999-07"
	parsed, err = time.Parse(layout, s)
	if err != nil {
		// Пробуем формат без микросекунд
		parsed, err = time.Parse("2006-01-02 15:04:05-07", s)
		if err != nil {
			return err
		}
	}
	*t = parsed
	return nil
})

var DataOptions = json.JoinOptions(
	json.WithUnmarshalers(timeUnmarshaler),
)
