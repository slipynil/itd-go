package posts

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

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
