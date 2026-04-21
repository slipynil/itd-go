package posts

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/slipynil/itd-go/internal/iterator"
	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// FeedIterator предоставляет интерфейс для постраничной загрузки постов.
type FeedIterator interface {
	// HasMore возвращает true, если есть ещё данные для загрузки.
	HasMore() bool
	// Next загружает и возвращает следующую страницу постов.
	// Параметры:
	//   - ctx: контекст для управления временем жизни запроса
	Next(ctx context.Context) ([]*types.Post, error)
}

// newFeedIterator создаёт итератор для получения ленты постов.
func newFeedIterator(ctx context.Context, s *Service, tab types.FeedTab, limit int) FeedIterator {
	fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.Post, iterator.PageToken, bool, error) {
		result, err := s.getFeed(ctx, tab, token.Cursor, limit)
		if err != nil {
			return nil, iterator.PageToken{}, false, err
		}
		next := iterator.PageToken{Cursor: result.Pagination.NextCursor}
		return result.Posts, next, result.Pagination.HasMore, nil
	}

	return iterator.New[*types.Post](fetch, iterator.PageToken{})
}

// newUserPostsIterator создаёт итератор для получения постов пользователя.
func newUserPostsIterator(ctx context.Context, s *Service, username string, limit int) FeedIterator {
	fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.Post, iterator.PageToken, bool, error) {
		result, err := s.getUserPosts(ctx, username, limit, token.Cursor)
		if err != nil {
			return nil, iterator.PageToken{}, false, err
		}
		next := iterator.PageToken{Cursor: result.Pagination.NextCursor}
		return result.Posts, next, result.Pagination.HasMore, nil
	}

	return iterator.New[*types.Post](fetch, iterator.PageToken{})
}

// getFeed возвращает api структуру с полями от получения постов
func (s *Service) getFeed(ctx context.Context, tab types.FeedTab, cursor string, limit int) (*FeedData, error) {

	path := fmt.Sprintf("/api/posts?limit=%d&tab=%v&cursor=%s", limit, tab, cursor)
	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result responseFeed
	if err := json.UnmarshalRead(resp.Body, &result, transport.DataOptions); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// getUserPosts возвращает api структуру с полями от получения постов пользователя
func (s *Service) getUserPosts(ctx context.Context, username string, limit int, cursor string) (*FeedData, error) {

	path := fmt.Sprintf("/api/posts/user/%s?limit=%d&sort=new&cursor=%s", username, limit, cursor)
	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result responseFeed

	if err := json.UnmarshalRead(resp.Body, &result, transport.DataOptions); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
