package types

import "time"

// Post представляет пост в социальной сети ITD.
// Содержит текстовый контент, метаданные, информацию об авторе и статистику взаимодействий.
type Post struct {
	// ID - уникальный идентификатор поста
	ID string `json:"id"`

	// Content - текстовое содержимое поста
	Content string `json:"content"`

	// Spans - массив элементов форматирования текста (жирный, курсив, ссылки)
	Spans []Span `json:"spans"`

	// LikesCount - количество лайков на посте
	LikesCount int `json:"likes_count"`

	// CommentsCount - количество комментариев к посту
	CommentsCount int `json:"comments_count"`

	// RepostsCount - количество репостов
	RepostsCount int `json:"reposts_count"`

	// ViewsCount - количество просмотров поста
	ViewsCount int `json:"views_count"`

	// IsLiked - true, если текущий пользователь лайкнул пост
	IsLiked bool `json:"is_liked"`

	// IsReposted - true, если текущий пользователь репостнул пост
	IsReposted bool `json:"is_reposted"`

	// IsPinned - true, если пост закреплён в профиле автора
	IsPinned bool `json:"is_pinned"`

	// CreatedAt - время создания поста
	CreatedAt time.Time `json:"created_at"`

	// Author - информация об авторе поста
	Author PostAuthor `json:"author"`

	// Attachments - массив вложений (изображения, видео, файлы)
	Attachments []Attachment `json:"attachments"`

	// Poll - опрос, прикреплённый к посту (может отсутствовать)
	Poll *Poll `json:"poll"`

	// OriginalPost - оригинальный пост, если текущий является репостом (может отсутствовать)
	OriginalPost *Post `json:"original_post"`

	// WallRecipient - получатель записи на стене (может отсутствовать)
	WallRecipient *PostAuthor `json:"wall_recipient"`
}

// Span представляет элемент форматирования текста в посте.
// Определяет тип форматирования (жирный, курсив, ссылка) и его позицию в тексте.
type Span struct {
	// Type - тип форматирования ("bold", "italic", "link", и т.д.)
	Type string `json:"type"`

	// Offset - позиция начала форматирования в тексте (в символах)
	Offset int `json:"offset"`

	// Length - длина форматируемого фрагмента (в символах)
	Length int `json:"length"`

	// URL - адрес ссылки (присутствует только для type: "link")
	URL *string `json:"url,omitempty"`
}

// PostAuthor представляет информацию об авторе поста или пользователе.
// Содержит базовые данные профиля и статус взаимоотношений с текущим пользователем.
type PostAuthor struct {
	// ID - уникальный идентификатор пользователя
	ID string `json:"id"`

	// Username - имя пользователя (логин)
	Username string `json:"username"`

	// DisplayName - отображаемое имя пользователя
	DisplayName string `json:"display_name"`

	// Avatar - URL аватара пользователя
	Avatar string `json:"avatar"`

	// IsVerified - true, если аккаунт верифицирован
	IsVerified bool `json:"is_verified"`

	// IsFollowing - true, если текущий пользователь подписан на этого пользователя
	IsFollowing bool `json:"is_following"`

	// IsFollowedBy - true, если этот пользователь подписан на текущего пользователя
	IsFollowedBy bool `json:"is_followed_by"`
}

// Attachment представляет файл, прикреплённый к посту.
// Может быть изображением, видео или другим типом файла.
type Attachment struct {
	// ID - уникальный идентификатор вложения
	ID string `json:"id"`

	// URL - прямая ссылка на файл
	URL string `json:"url"`

	// Type - тип вложения ("image", "video", "file", и т.д.)
	Type string `json:"type"`

	// Size - размер файла в байтах
	Size int64 `json:"size"`

	// Width - ширина изображения в пикселях (присутствует только для type: "image")
	Width *int `json:"width,omitempty"`

	// Height - высота изображения в пикселях (присутствует только для type: "image")
	Height *int `json:"height,omitempty"`
}

// Poll представляет опрос, прикреплённый к посту.
// Содержит вопрос, варианты ответов и информацию о голосовании.
type Poll struct {
	// ID - уникальный идентификатор опроса
	ID string `json:"id"`

	// PostID - идентификатор поста, к которому прикреплён опрос
	PostID string `json:"post_id"`

	// Question - текст вопроса опроса
	Question string `json:"question"`

	// Options - массив вариантов ответа
	Options []PollOption `json:"options"`

	// TotalVotes - общее количество голосов в опросе
	TotalVotes int `json:"total_votes"`

	// HasVoted - true, если текущий пользователь проголосовал
	HasVoted bool `json:"has_voted"`

	// VotedOptionIDs - массив ID вариантов, за которые проголосовал текущий пользователь
	VotedOptionIDs []string `json:"voted_option_ids"`

	// MultipleChoice - true, если можно выбрать несколько вариантов ответа
	MultipleChoice bool `json:"multiple_choice"`
}

// PollOption представляет один вариант ответа в опросе.
type PollOption struct {
	// ID - уникальный идентификатор варианта ответа
	ID string `json:"id"`

	// Text - текст варианта ответа
	Text string `json:"text"`

	// Votes - количество голосов за этот вариант
	Votes int `json:"votes"`
}
