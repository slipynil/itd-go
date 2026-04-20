package posts

import (
	"context"

	"github.com/slipynil/itd-go/internal/pkg/iterator"
	"github.com/slipynil/itd-go/types"
)

// newFeedIterator создаёт итератор для получения ленты постов.
func newFeedIterator(s *Service, ctx context.Context, tab types.FeedTab, limit int) types.FeedIterator {
	fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.Post, iterator.PageToken, bool, error) {
		result, err := s.getFeed(ctx, tab, token.Cursor, limit)
		if err != nil {
			return nil, iterator.PageToken{}, false, err
		}
		next := iterator.PageToken{Cursor: result.Pagination.NextCursor}
		return result.Posts, next, result.Pagination.HasMore, nil
	}

	return iterator.New[*types.Post](ctx, fetch, iterator.PageToken{})
}

// newUserPostsIterator создаёт итератор для получения постов пользователя.
func newUserPostsIterator(s *Service, ctx context.Context, username string, limit int) types.FeedIterator {
	fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.Post, iterator.PageToken, bool, error) {
		result, err := s.getUserPosts(ctx, username, limit, token.Cursor)
		if err != nil {
			return nil, iterator.PageToken{}, false, err
		}
		next := iterator.PageToken{Cursor: result.Pagination.NextCursor}
		return result.Posts, next, result.Pagination.HasMore, nil
	}

	return iterator.New[*types.Post](ctx, fetch, iterator.PageToken{})
}
