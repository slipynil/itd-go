package posts

import (
	"context"

	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

type Posts struct {
	transport *transport.Client
}

func New(t *transport.Client) *Posts {
	return &Posts{transport: t}
}

func (p *Posts) NewFeed(ctx context.Context, limit int, sort string) types.FeedIterator {
	return newFeedIterator(p, ctx, limit, sort)
}
