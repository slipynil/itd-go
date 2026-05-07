package search

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

type Service struct {
	transport *transport.Client
}

func New(transport *transport.Client) *Service {
	return &Service{transport: transport}
}

// NewHashtagPosts создаёт итератор для получения постов по хештегу.
// Параметры:
//   - hashtag: название хештега (без #)
//   - limit: количество постов на страницу (от 1 до 50)
//
// Возвращает Iterator для постраничной загрузки постов или ошибку при невалидном limit.
func (s *Service) NewHashtagPosts(hashtag string, limit int) (Iterator, error) {
	if limit < 1 || limit > 50 {
		return nil, fmt.Errorf("limit must be between 1 and 50")
	}
	return newHashtagPosts(s, hashtag, limit), nil
}

func (s *Service) TopHashtags(ctx context.Context, limit int) ([]types.Hashtag, error) {
	if limit < 1 || limit > 50 {
		return nil, fmt.Errorf("limit must be between 1 and 50")
	}
	path := fmt.Sprintf("/api/hashtags/trending?limit=%d", limit)
	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result topHashtagsResponse
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}
	return result.Data.Hashtag, nil
}

func (s *Service) Top10Clans(ctx context.Context) ([]types.Clans, error) {
	path := "/api/users/stats/top-clans"
	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result topClansResponse
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}
	return result.Clans, nil
}

// getHashtagFeed возвращает посты по хештегу с пагинацией.
// Используется внутри итератора для загрузки страниц.
func (s *Service) getHashtagFeed(ctx context.Context, hashtag string, cursor string, limit int) (*FeedData, error) {
	path := fmt.Sprintf("/api/hashtags/%s/posts?limit=%d", hashtag, limit)
	if cursor != "" {
		path = fmt.Sprintf("%s&cursor=%s", path, cursor)
	}

	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result hashtagFeedResponse
	if err := json.UnmarshalRead(resp.Body, &result, transport.DataOptions); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (s *Service) Query(ctx context.Context, query string) (*types.SearchResult, error) {
	path := fmt.Sprintf("/api/search/?q=%s&userLimit=20&hashtagLimit=20", query)
	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result searchResponse
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
