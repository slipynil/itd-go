package types

import "time"

// Comment представляет комментарий к посту в ITD.
// Содержит текстовый контент, информацию об авторе и статистику взаимодействий.
type Comment struct {
	// ID - уникальный идентификатор комментария
	ID string `json:"id"`

	// Content - текстовое содержимое комментария
	Content string `json:"content"`

	// Author - информация об авторе комментария
	Author PostAuthor `json:"author"`

	// LikesCount - количество лайков на комментарии
	LikesCount int `json:"likesCount"`

	// RepliesCount - количество ответов на комментарий
	RepliesCount int `json:"repliesCount"`

	// IsLiked - true, если текущий пользователь лайкнул комментарий
	IsLiked bool `json:"isLiked"`

	// CreatedAt - время создания комментария
	CreatedAt time.Time `json:"createdAt"`

	// Attachments - массив вложений (изображения, файлы)
	Attachments []*CommentAttachment `json:"attachments"`

	// Replies - массив ответов на комментарий (может быть nil)
	Replies []*Comment `json:"replies,omitempty"`

	// ReplyTo - информация о комментарии, на который отвечает данный комментарий (nil, если это не ответ)
	ReplyTo *ReplyTo `json:"replyTo,omitempty"`
}

// CommentUpdate представляет результат обновления комментария.
type CommentUpdate struct {
	// ID - уникальный идентификатор комментария
	ID string `json:"id"`

	// Content - обновлённое текстовое содержимое
	Content string `json:"content"`

	// UpdatedAt - время обновления комментария
	UpdatedAt time.Time `json:"editedAt"`
}

// Author представляет краткую информацию об авторе комментария.
// Содержит базовые данные профиля пользователя.
type Author struct {
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

// CommentAttachment представляет файл, прикреплённый к комментарию.
// Может быть изображением, видео или другим типом файла.
type CommentAttachment struct {
	// ID - уникальный идентификатор вложения
	ID string `json:"id"`

	// Type - тип вложения ("image", "video", "file" и т.д.)
	Type string `json:"type"`

	// URL - публичная ссылка на файл
	URL string `json:"url"`

	// Filename - имя файла
	Filename string `json:"filename"`

	// MimeType - MIME-тип файла (например, "image/png", "video/mp4")
	MimeType string `json:"mimeType"`

	// Size - размер файла в байтах
	Size int `json:"size"`

	// Width - ширина в пикселях (присутствует только для изображений и видео)
	Width *int `json:"width,omitempty"`

	// Height - высота в пикселях (присутствует только для изображений и видео)
	Height *int `json:"height,omitempty"`

	// Duration - длительность в секундах (присутствует только для видео и аудио)
	Duration *int `json:"duration,omitempty"`

	// Order - порядковый номер вложения в комментарии
	Order int `json:"order"`
}

// ReplyTo представляет информацию о пользователе, на чей комментарий был дан ответ.
// Содержит минимальные данные для отображения ссылки на оригинальный комментарий.
type ReplyTo struct {
	// ID - уникальный идентификатор пользователя
	ID string `json:"id"`

	// Username - уникальный логин пользователя
	Username string `json:"username"`

	// DisplayName - отображаемое имя пользователя
	DisplayName string `json:"displayName"`
}

type CreateComment struct {
	// ID - уникальный идентификатор комментария
	ID string `json:"id"`

	// Content - текстовое содержимое комментария
	Content string `json:"content"`

	// Author - информация об авторе комментария
	Author PostAuthor `json:"author"`

	// CreatedAt - время создания комментария
	CreatedAt time.Time `json:"createdAt"`

	// Attachments - массив вложений (изображения, файлы)
	Attachments []*CommentAttachment `json:"attachments"`
}
