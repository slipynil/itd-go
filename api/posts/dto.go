package posts

import "github.com/slipynil/itd-go/types"

// createPostRequest представляет запрос на создание поста.
type createPostRequest struct {
	Content       string             `json:"content,omitempty"`
	Spans         []types.Span       `json:"spans,omitempty"`
	AttachmentIDs []string           `json:"attachmentIds,omitempty"`
	Poll          *types.PollRequest `json:"poll,omitempty"`
}

// repostRequest представляет запрос на создание репоста.
type repostRequest struct {
	Content string `json:"content,omitempty"`
}

// voteRequest представляет запрос на голосование в опросе.
type voteRequest struct {
	OptionIds []string `json:"optionIds,omitempty"`
}

// ResponseFeed представляет ответ API при получении ленты постов.
type responseFeed struct {
	// Data - данные ленты с постами и информацией о пагинации
	Data FeedData `json:"data"`
}

// ResponsePost представляет ответ API при получении одного поста.
type responsePost struct {
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
