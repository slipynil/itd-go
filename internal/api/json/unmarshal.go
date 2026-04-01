package json

import (
	"io"
	"strings"
	"time"

	"github.com/go-json-experiment/json"
)

var timeUnmarshaler = json.UnmarshalFunc(func(b []byte, t *time.Time) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	layout := "2006-01-02 15:04:05.999999-07"
	parsed, err := time.Parse(layout, s)
	if err != nil {
		parsed, err = time.Parse("2006-01-02 15:04:05-07", s)
		if err != nil {
			return err
		}
	}
	*t = parsed
	return nil
})

var defaultOptions = json.JoinOptions(
	json.WithUnmarshalers(timeUnmarshaler),
	json.MatchCaseInsensitiveNames(true),
)

// Unmarshal - публичная функция-обертка для удобного парсинга
func Unmarshal(data io.Reader, v any) error {
	return json.UnmarshalRead(data, v, defaultOptions)
}
