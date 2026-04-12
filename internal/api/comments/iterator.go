package comments

import (
	"context"

	"github.com/slipynil/itd-go/internal/pkg/iterator"
	"github.com/slipynil/itd-go/types"
)

func commentListIterator(client *Comments, ctx context.Context, postID string, limit int) types.CommentIterator {
	fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.Comment, iterator.PageToken, bool, error) {
		result, err := client.getCommentList(ctx, postID, token.Cursor, limit)
		if err != nil {
			return nil, iterator.PageToken{}, false, err
		}
		next := iterator.PageToken{Cursor: result.Data.NextCursor}
		return result.Data.Comments, next, result.Data.HasMore, nil
	}

	return iterator.New[*types.Comment](ctx, fetch, iterator.PageToken{})
}
