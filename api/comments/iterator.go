package comments

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/slipynil/itd-go/internal/iterator"
	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// CommentIterator предоставляет интерфейс для постраничной загрузки комментариев.
type CommentIterator interface {
	// HasMore возвращает true, если есть ещё данные для загрузки.
	HasMore() bool
	// Next загружает и возвращает следующую страницу комментариев.
	// Параметры:
	//   - ctx: контекст для управления временем жизни запроса
	Next(ctx context.Context) ([]*types.Comment, error)
}

func commentListIterator(s *Service, postID string, limit int) CommentIterator {
	fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.Comment, iterator.PageToken, bool, error) {
		result, err := s.getCommentList(ctx, postID, token.Cursor, limit)
		if err != nil {
			return nil, iterator.PageToken{}, false, err
		}
		next := iterator.PageToken{Cursor: result.Data.NextCursor}
		return result.Data.Comments, next, result.Data.HasMore, nil
	}

	return iterator.New[*types.Comment](fetch, iterator.PageToken{})
}

// getCommentList получает сырую json стурктуру с комментариями и информацией о пагинации.
func (s *Service) getCommentList(ctx context.Context, postID, cursor string, limit int) (*commentsResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/comments?limit=%d&sort=popular", postID, limit)

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
