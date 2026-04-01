package posts

import (
	"context"

	"github.com/slipynil/itd-go/types"
)

type feedIterator struct {
	client  *Posts
	ctx     context.Context
	limit   int
	sort    string
	cursor  string
	hasMore bool
}

func newFeedIterator(client *Posts, ctx context.Context, limit int, sort string) types.FeedIterator {
	return &feedIterator{
		client:  client,
		ctx:     ctx,
		limit:   limit,
		sort:    sort,
		hasMore: true,
	}
}

func (f *feedIterator) HasMore() bool {
	return f.hasMore
}

func (f *feedIterator) Next() ([]types.Post, error) {
	result, err := f.client.getFeed(f.ctx, f.limit, f.cursor, f.sort)
	if err != nil {
		return nil, err
	}
	f.cursor = result.Pagination.NextCursor
	f.hasMore = result.Pagination.HasMore

	return result.Posts, nil
}
