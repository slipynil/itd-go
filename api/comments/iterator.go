package comments

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/slipynil/itd-go/internal/iterator"
	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// Iterator предоставляет интерфейс для постраничной загрузки комментариев.
type Iterator interface {
	// HasMore возвращает true, если есть ещё данные для загрузки.
	HasMore() bool
	// Next загружает и возвращает следующую страницу комментариев.
	// Параметры:
	//   - ctx: контекст для управления временем жизни запроса
	Next(ctx context.Context) ([]*types.Comment, error)
}

// newPostComments создаёт итератор для получения комментариев к посту.
// Параметры:
//   - s: сервис для работы с API комментариев
//   - postID: идентификатор поста
//   - limit: количество комментариев на страницу
func newPostComments(s *Service, postID string, limit int) Iterator {
	fetch := func(ctx context.Context, token *iterator.PageToken) ([]*types.Comment, *iterator.PageToken, bool, error) {
		cursor := ""
		if token != nil {
			cursor = token.Cursor
		}

		result, err := s.getCommentList(ctx, postID, cursor, limit)
		if err != nil {
			return nil, nil, false, err
		}

		var next *iterator.PageToken
		if result.Data.HasMore {
			next = &iterator.PageToken{Cursor: result.Data.NextCursor}
		}

		return result.Data.Comments, next, result.Data.HasMore, nil
	}

	return iterator.New[*types.Comment](fetch, nil)
}

// getCommentList получает комментарии к посту с пагинацией.
// Используется внутри итератора для загрузки страниц.
func (s *Service) getCommentList(ctx context.Context, postID, cursor string, limit int) (*commentsResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/comments?limit=%d&sort=popular", postID, limit)
	if cursor != "" {
		path = fmt.Sprintf("%s&cursor=%s", path, cursor)
	}

	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result commentsResponse
	if err := json.UnmarshalRead(resp.Body, &result, transport.DataOptions); err != nil {
		return nil, err
	}

	return &result, nil
}
