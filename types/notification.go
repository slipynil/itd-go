package types

import "time"

// Notification представляет уведомление из API ITD.
// Используется методами ListUnread, NewIterator и другими REST endpoints.
type Notification struct {
	// ID - уникальный идентификатор уведомления
	ID string `json:"id"`
	// Type - тип уведомления (like, comment, follow и т.д.)
	Type string `json:"type"`
	// TargetType - тип целевого объекта (post, comment и т.д.)
	TargetType string `json:"targetType"`
	// TargetID - идентификатор целевого объекта
	TargetID string `json:"targetId"`
	// Preview - текстовое превью уведомления
	Preview string `json:"preview"`
	// ReadAt - время прочтения уведомления (может быть zero value если не прочитано)
	ReadAt time.Time `json:"readAt"`
	// CreatedAt - время создания уведомления
	CreatedAt time.Time `json:"createdAt"`
	// Actor - информация о пользователе, вызвавшем уведомление
	Actor Actor `json:"actor"`
	// Read - флаг прочтения уведомления
	Read bool `json:"read"`
}

// StreamNotification представляет уведомление из SSE стрима.
// Используется методом Stream для получения уведомлений в реальном времени.
// Отличается от Notification наличием дополнительных полей и nullable ReadAt.
type StreamNotification struct {
	// ID - уникальный идентификатор уведомления
	ID string `json:"id"`
	// Type - тип уведомления (like, comment, follow и т.д.)
	Type string `json:"type"`
	// TargetType - тип целевого объекта (post, comment и т.д.)
	TargetType string `json:"targetType"`
	// TargetID - идентификатор целевого объекта
	TargetID string `json:"targetId"`
	// Preview - текстовое превью уведомления
	Preview string `json:"preview"`
	// ReadAt - время прочтения уведомления (nil если не прочитано)
	ReadAt *time.Time `json:"readAt"`
	// CreatedAt - время создания уведомления
	CreatedAt time.Time `json:"createdAt"`
	// UserID - идентификатор пользователя-получателя уведомления
	UserID string `json:"userId"`
	// Actor - информация о пользователе, вызвавшем уведомление
	Actor Actor `json:"actor"`
	// Read - флаг прочтения уведомления
	Read bool `json:"read"`
	// Sound - флаг необходимости воспроизведения звука
	Sound bool `json:"sound"`
}

// Actor представляет информацию о пользователе, вызвавшем уведомление.
// Содержит базовую информацию о профиле и статус подписки.
type Actor struct {
	// ID - уникальный идентификатор пользователя
	ID string `json:"id"`
	// DisplayName - отображаемое имя пользователя
	DisplayName string `json:"displayName"`
	// Username - уникальное имя пользователя (без @)
	Username string `json:"username"`
	// Avatar - URL аватара пользователя
	Avatar string `json:"avatar"`
	// IsFollowing - подписан ли текущий пользователь на актора
	IsFollowing bool `json:"isFollowing"`
	// IsFollowedBy - подписан ли актор на текущего пользователя
	IsFollowedBy bool `json:"isFollowedBy"`
}
