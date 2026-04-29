package search

import (
	"context"

	"github.com/slipynil/itd-go/internal/iterator"
	"github.com/slipynil/itd-go/types"
)

// Iterator предоставляет интерфейс для постраничной загрузки постов по хештегу.
type Iterator interface {
	// HasMore возвращает true, если есть ещё данные для загрузки.
	HasMore() bool
	// Next загружает и возвращает следующую страницу постов.
	// Параметры:
	//   - ctx: контекст для управления временем жизни запроса
	Next(ctx context.Context) ([]*types.Post, error)
}

// newHashtagPosts создаёт итератор для получения постов по хештегу.
// Параметры:
//   - s: сервис для работы с API поиска
//   - hashtag: название хештега (без #)
//   - limit: количество постов на страницу
func newHashtagPosts(s *Service, hashtag string, limit int) Iterator {
	fetch := func(ctx context.Context, token *iterator.PageToken) ([]*types.Post, *iterator.PageToken, bool, error) {
		cursor := ""
		if token != nil {
			cursor = token.Cursor
		}

		result, err := s.getHashtagFeed(ctx, hashtag, cursor, limit)
		if err != nil {
			return nil, nil, false, err
		}

		var next *iterator.PageToken
		if result.Pagination.HasMore {
			next = &iterator.PageToken{Cursor: result.Pagination.NextCursor}
		}

		return result.Posts, next, result.Pagination.HasMore, nil
	}

	return iterator.New[*types.Post](fetch, nil)
}
