package posts

import "github.com/slipynil/itd-go/types"

// ResponseFeed представляет ответ API при получении ленты постов.
type ResponseFeed struct {
	// Data - данные ленты с постами и информацией о пагинации
	Data FeedData `json:"data"`
}

// ResponsePost представляет ответ API при получении одного поста.
type ResponsePost struct {
	// Data - данные поста
	Data types.Post `json:"data"`
}

// FeedData содержит массив постов и информацию о пагинации.
type FeedData struct {
	// Posts - массив постов на текущей странице
	Posts []*types.Post `json:"posts"`

	// Pagination - информация о пагинации
	Pagination Pagination `json:"pagination"`
}

// Pagination содержит информацию о пагинации для ленты постов.
type Pagination struct {
	// Limit - количество постов на странице
	Limit int `json:"limit"`

	// NextCursor - курсор для получения следующей страницы
	NextCursor string `json:"nextCursor"`

	// HasMore - true, если есть ещё посты для загрузки
	HasMore bool `json:"hasMore"`
}
