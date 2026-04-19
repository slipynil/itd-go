package types

// AuthorInfo представляет краткую информацию об авторе контента.
// Используется в постах, комментариях и других сущностях.
type AuthorInfo struct {
	// ID - уникальный идентификатор пользователя
	ID string `json:"id"`

	// Username - уникальный логин пользователя
	Username string `json:"username"`

	// DisplayName - отображаемое имя пользователя
	DisplayName string `json:"displayName"`

	// Avatar - аватар пользователя: эмодзи или URL изображения
	Avatar string `json:"avatar"`

	// Verified - true, если аккаунт верифицирован
	Verified bool `json:"verified"`

	// Pin - специальный значок пользователя (nil, если не установлен)
	Pin *Pin `json:"pin,omitempty"`

	// HasNuksta - true, если пользователь имеет премиум подписку
	HasNuksta bool `json:"hasNuksta"`
}
