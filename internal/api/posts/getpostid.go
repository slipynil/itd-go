package posts

import (
	"context"
	"fmt"

	"github.com/slipynil/itd-go/internal/api/dto"
	"github.com/slipynil/itd-go/internal/api/json"
	"github.com/slipynil/itd-go/types"
)

func (p *Posts) GetPostByID(ctx context.Context, id int) (*types.Post, error) {
	path := fmt.Sprintf("api/posts/%s", id)
	req, err := p.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var post types.Post
	if err := json.Unmarshal(resp.Body, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *Posts) getFeed(ctx context.Context, limit int, cursor, sort string) (*dto.FeedData, error) {
	path := fmt.Sprintf("/api/posts?limit=%d&tab=%s&cursor=%s", limit, sort, cursor)
	req, err := p.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.transport.Do(req)
	if err != nil {
		return nil, err
	}
	var feed dto.FeedResponse
	if err := json.Unmarshal(resp.Body, &feed); err != nil {
		return nil, err
	}
	return &feed.Data, nil
}
