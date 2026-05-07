package posts

import (
	"context"

	"github.com/slipynil/itd-go/internal/iterator"
	"github.com/slipynil/itd-go/types"
)

// Iterator предоставляет интерфейс для постраничной загрузки постов.
type Iterator interface {
	// HasMore возвращает true, если есть ещё данные для загрузки.
	HasMore() bool
	// Next загружает и возвращает следующую страницу постов.
	// Параметры:
	//   - ctx: контекст для управления временем жизни запроса
	Next(ctx context.Context) ([]*types.Post, error)
}

// newFeed создаёт итератор для получения ленты постов.
// Параметры:
//   - s: сервис для работы с API постов
//   - tab: тип ленты (popular, clan, following)
//   - limit: количество постов на страницу
func newFeed(s *Service, tab types.FeedTab, limit int) Iterator {
	fetch := func(ctx context.Context, token *iterator.PageToken) ([]*types.Post, *iterator.PageToken, bool, error) {
		cursor := ""
		if token != nil {
			cursor = token.Cursor
		}

		result, err := s.getFeed(ctx, tab, cursor, limit)
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

// newUserPosts создаёт итератор для получения постов пользователя.
// Параметры:
//   - s: сервис для работы с API постов
//   - username: имя пользователя
//   - limit: количество постов на страницу
func newUserPosts(s *Service, username string, limit int) Iterator {
	fetch := func(ctx context.Context, token *iterator.PageToken) ([]*types.Post, *iterator.PageToken, bool, error) {
		cursor := ""
		if token != nil {
			cursor = token.Cursor
		}

		result, err := s.getUserPosts(ctx, username, limit, cursor)
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
