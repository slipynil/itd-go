package types

import "time"

// CreatedPostBase содержит общие поля для всех типов созданных постов.
// Используется как встраиваемая структура в CreatedPost, CreatedPostWithPoll и CreatedPostWithRepost.
type CreatedPostBase struct {
	// ID - уникальный идентификатор поста
	ID string `json:"id"`

	// Content - текстовое содержимое поста (может содержать эмодзи)
	Content string `json:"content"`

	// CreatedAt - дата и время создания поста
	CreatedAt time.Time `json:"createdAt"`

	// Attachments - массив вложений (изображения, видео, файлы)
	Attachments []Attachment `json:"attachments"`
}

// CreatedPost представляет созданный пост в социальной сети ITD.
// Содержит id, текстовый контент и информацию о времени создания.
type CreatedPost struct {
	CreatedPostBase
}

// CreatedPostWithPoll представляет созданный пост с опросом в социальной сети ITD.
// Содержит id, текстовый контент, информацию о времени создания и опрос.
type CreatedPostWithPoll struct {
	CreatedPostBase

	// Poll - опрос, прикреплённый к посту
	Poll *Poll `json:"poll"`
}

// CreatedPostWithRepost представляет созданный пост с репостом в социальной сети ITD.
// Содержит оригинальный пост
type CreatedPostWithRepost struct {
	CreatedPostBase

	// OriginalPost - оригинальный пост, если текущий является репостом (nil, если не репост)
	OriginalPost *Post `json:"originalPost"`
}

// Post представляет пост в социальной сети ITD.
// Содержит текстовый контент, метаданные, информацию об авторе и статистику взаимодействий.
type Post struct {
	// ID - уникальный идентификатор поста
	ID string `json:"id"`

	// Content - текстовое содержимое поста (может содержать эмодзи)
	Content string `json:"content"`

	// Spans - массив элементов форматирования текста (жирный, курсив, ссылки и т.д.)
	Spans []Span `json:"spans"`

	// LikesCount - количество лайков на посте
	LikesCount int `json:"likesCount"`

	// CommentsCount - количество комментариев к посту
	CommentsCount int `json:"commentsCount"`

	// RepostsCount - количество репостов
	RepostsCount int `json:"repostsCount"`

	// ViewsCount - количество просмотров поста
	ViewsCount int `json:"viewsCount"`

	// IsLiked - true, если текущий пользователь лайкнул пост
	IsLiked bool `json:"isLiked"`

	// IsReposted - true, если текущий пользователь репостнул пост
	IsReposted bool `json:"isReposted"`

	// IsOwner - true, если текущий пользователь является владельцем поста
	IsOwner bool `json:"isOwner"`

	// IsViewed - true, если текущий пользователь уже просмотрел пост
	IsViewed bool `json:"isViewed"`

	// IsDeleted - true, если пост был удалён
	IsDeleted bool `json:"isDeleted"`

	// CreatedAt - дата и время создания поста
	CreatedAt time.Time `json:"createdAt"`

	// Author - информация об авторе поста
	Author AuthorInfo `json:"author"`

	// Attachments - массив вложений (изображения, видео, файлы)
	Attachments []Attachment `json:"attachments"`

	// Poll - опрос, прикреплённый к посту (nil, если опроса нет)
	Poll *Poll `json:"poll"`

	// OriginalPost - оригинальный пост, если текущий является репостом (nil, если не репост)
	OriginalPost *Post `json:"originalPost"`

	// WallRecipientID - ID получателя записи на стене (nil, если пост не на чужой стене)
	WallRecipientID *string `json:"wallRecipientId"`

	// EditedAt - дата и время последнего редактирования (nil, если пост не редактировался)
	EditedAt *time.Time `json:"editedAt"`

	// DominantEmoji - доминирующий эмодзи поста (пустая строка, если отсутствует)
	DominantEmoji string `json:"dominantEmoji"`
}

// Span представляет элемент форматирования текста в посте.
// Определяет тип форматирования (жирный, курсив, ссылка) и его позицию в тексте.
type Span struct {
	// Type - тип форматирования ("bold", "italic", "link" и т.д.)
	Type string `json:"type"`

	// Length - длина форматируемого фрагмента (в символах)
	Length int `json:"length"`

	// Offset - позиция начала форматирования в тексте (в символах)
	Offset int `json:"offset"`

	// Эти поля специфичны для разных типов Span
	Username string `json:"username,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

// Attachment представляет файл, прикреплённый к посту.
// Может быть изображением, видео или другим типом файла.
type Attachment struct {
	// ID - уникальный идентификатор вложения
	ID string `json:"id"`

	// URL - публичная ссылка на файл
	URL string `json:"url"`

	// Type - тип вложения ("image", "video", "file" и т.д.)
	Type string `json:"type"`

	// Size - размер файла в байтах (может быть 0, если сервер не возвращает значение)
	Size int64 `json:"size"`

	// Width - ширина в пикселях (присутствует только при type: "image")
	Width *int `json:"width,omitempty"`

	// Height - высота в пикселях (присутствует только при type: "image")
	Height *int `json:"height,omitempty"`
}

// Poll представляет опрос, прикреплённый к посту.
// Содержит вопрос, варианты ответов и информацию о голосовании.
type Poll struct {
	// ID - уникальный идентификатор опроса
	ID string `json:"id"`

	// PostID - идентификатор поста, к которому прикреплён опрос
	PostID string `json:"postId"`

	// Question - текст вопроса опроса
	Question string `json:"question"`

	// Options - массив вариантов ответа
	Options []PollOptionResponse `json:"options"`

	// TotalVotes - суммарное количество голосов по всем вариантам
	TotalVotes int `json:"totalVotes"`

	// HasVoted - true, если текущий пользователь уже проголосовал
	HasVoted bool `json:"hasVoted"`

	// VotedOptionIDs - ID вариантов, за которые проголосовал текущий пользователь
	VotedOptionIDs []string `json:"votedOptionIds"`

	// MultipleChoice - true, если допускается выбор нескольких вариантов
	MultipleChoice bool `json:"multipleChoice"`
}

// PollOptionResponse представляет один вариант ответа в опросе.
type PollOptionResponse struct {
	// ID - уникальный идентификатор варианта ответа
	ID string `json:"id"`

	// Text - текст варианта ответа
	Text string `json:"text"`

	// Votes - количество голосов за этот вариант
	Votes int `json:"votesCount"`
}

// PollRequest представляет запрос на создание опроса.
type PollRequest struct {
	// Question - текст вопроса опроса
	Question string `json:"question"`

	// Options - массив вариантов ответа
	Options []PollOptionRequest `json:"options"`

	// MultipleChoice - true, если допускается выбор нескольких вариантов
	MultipleChoice bool `json:"multipleChoice"`
}

// PollOptionRequest представляет один вариант ответа в запросе на создание опроса.
type PollOptionRequest struct {
	// Text - текст варианта ответа
	Text string `json:"text"`
}

// LikesCountResponse представляет ответ API с количеством лайков.
// Возвращается методами Like и Unlike.
type LikesCountResponse struct {
	// LikesCount - текущее количество лайков после операции
	LikesCount int `json:"likesCount"`
}

// PostViewResponse представляет ответ API при отметке просмотра поста.
type PostViewResponse struct {
	// Viewed - true, если просмотр успешно зарегистрирован
	Viewed bool `json:"viewed"`
}
