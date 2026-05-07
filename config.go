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

	// MaxRetries - максимальное количество повторных попыток при ошибке 429 (rate limiting).
	// По умолчанию 3. Установите 0 для отключения retry логики.
	MaxRetries int

	// RetryDelay - начальная задержка перед первой повторной попыткой.
	// Используется exponential backoff: delay, delay*2, delay*4, и т.д.
	// По умолчанию 1 секунда.
	RetryDelay time.Duration
}
