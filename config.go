package itdgo

import (
	"time"
)

// Config содержит параметры конфигурации для ITD SDK клиента.
// RefreshToken является обязательным параметром, остальные поля опциональны.
type Config struct {
	// RefreshToken - обязательный параметр для аутентификации.
	// Получить можно из cookies браузера (cookie с именем "refresh_token").
	RefreshToken string

	// UserAgent - User-Agent заголовок для HTTP запросов.
	// Если не указан, будет использован User-Agent по умолчанию.
	UserAgent string

	// Timeout - таймаут для HTTP запросов.
	// Если не указан, используется значение по умолчанию (30 секунд).
	Timeout time.Duration

	// WithoutBanner - если true, баннер SDK не будет выведен при инициализации.
	// По умолчанию false (баннер выводится).
	WithoutBanner bool
}
