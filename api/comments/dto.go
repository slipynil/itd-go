package comments

import "github.com/slipynil/itd-go/types"

// createCommentRequest представляет запрос на создание комментария.
type createCommentRequest struct {
	Content       string   `json:"content,omitempty"`
	AttachmentIDs []string `json:"attachmentIds,omitempty"`
}

// createReplyRequest представляет запрос на создание ответа на комментарий.
type createReplyRequest struct {
	ReplyToUserId string   `json:"replyToUserId,omitempty"`
	Content       string   `json:"content,omitempty"`
	AttachmentIDs []string `json:"attachmentIds,omitempty"`
}

// updateCommentRequest представляет запрос на обновление комментария.
type updateCommentRequest struct {
	Content string `json:"content,omitempty"`
}

// Response представляет ответ API при получении списка комментариев.
type commentsResponse struct {
	// Data - данные с комментариями и информацией о пагинации
	Data CommentsData `json:"data"`
}

type repliesResponse struct {
	Data RepliesData `json:"data"`
}

type RepliesData struct {
	Replies    []*types.Comment `json:"replies"`
	Pagination Pagination       `json:"pagination"`
}

type Pagination struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasMore bool `json:"hasMore"`
}

// Data содержит массив комментариев и информацию о пагинации.
type CommentsData struct {
	// Comments - массив комментариев на текущей странице
	Comments []*types.Comment `json:"comments"`

	// Total - общее количество комментариев
	Total int `json:"total"`

	// HasMore - true, если есть ещё комментарии для загрузки
	HasMore bool `json:"hasMore"`

	// NextCursor - курсор для получения следующей страницы
	NextCursor string `json:"nextCursor"`
}
