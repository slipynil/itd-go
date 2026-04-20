package posts

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-json-experiment/json"

	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// Posts предоставляет методы для работы с постами ITD API.
type Service struct {
	transport *transport.Client
}

// New создаёт новый экземпляр клиента для работы с постами.
func New(t *transport.Client) *Service {
	return &Service{transport: t}
}

// NewFeed возвращает итератор для получения постов.
func (s *Service) NewFeed(ctx context.Context, tab types.FeedTab, limit int) types.FeedIterator {
	return newFeedIterator(s, ctx, tab, limit)
}

// NewUserPosts возвращает итератор для получения постов пользователя.
func (s *Service) NewUserPosts(ctx context.Context, username string, limit int) types.FeedIterator {
	return newUserPostsIterator(s, ctx, username, limit)
}

// Get получает пост по его ID.
func (s *Service) Get(ctx context.Context, postID string) (*types.Post, error) {
	path := fmt.Sprintf("/api/posts/%s", postID)
	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.Post `json:"data"`
	}

	err = json.UnmarshalRead(resp.Body, &result, transport.DataOptions)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// Create создаёт новый пост.
func (s *Service) Create(ctx context.Context, content string, filePaths ...string) (*types.Post, error) {
	if strings.TrimSpace(content) == "" && len(filePaths) == 0 {
		return nil, fmt.Errorf("content or files required")
	}

	attachmentIDs, err := s.transport.UploadFiles(ctx, filePaths...)
	if err != nil {
		return nil, err
	}

	payload := createPostRequest{
		Content:       content,
		AttachmentIDs: attachmentIDs,
	}

	req, err := s.transport.NewRequest(ctx, "POST", "/api/posts", payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.Post

	err = json.UnmarshalRead(resp.Body, &result, transport.DataOptions)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateWithPoll создаёт новый пост с опросом.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - content: текстовое содержимое поста
//   - poll: структура опроса с вопросом и вариантами ответов
//   - filePaths: пути к файлам для загрузки и прикрепления к посту
//
// Возвращает созданный пост с опросом или ошибку при проблемах с сетью/API.
func (s *Service) CreateWithPoll(ctx context.Context, content string, poll *types.PollRequest, filePaths ...string) (*types.Post, error) {
	if poll == nil {
		return nil, fmt.Errorf("poll cannot be nil")
	}
	if len(poll.Options) < 2 {
		return nil, fmt.Errorf("poll must have at least 2 options")
	}

	attachmentIDs, err := s.transport.UploadFiles(ctx, filePaths...)
	if err != nil {
		return nil, err
	}

	payload := createPostRequest{
		Content:       content,
		AttachmentIDs: attachmentIDs,
		Poll:          poll,
	}

	req, err := s.transport.NewRequest(ctx, "POST", "/api/posts", payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.Post

	err = json.UnmarshalRead(resp.Body, &result, transport.DataOptions)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete удаляет пост по его ID.
func (s *Service) Delete(ctx context.Context, postID string) error {
	path := fmt.Sprintf("/api/posts/%s", postID)
	req, err := s.transport.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Like ставит лайк на пост.
func (s *Service) Like(ctx context.Context, postID string) (*types.LikesCountResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/like", postID)
	req, err := s.transport.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.LikesCountResponse

	err = json.UnmarshalRead(resp.Body, &result, transport.DataOptions)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Unlike убирает лайк с поста.
func (s *Service) Unlike(ctx context.Context, postID string) (*types.LikesCountResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/like", postID)
	req, err := s.transport.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.LikesCountResponse

	err = json.UnmarshalRead(resp.Body, &result, transport.DataOptions)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Repost создаёт репост существующего поста.
func (s *Service) Repost(ctx context.Context, postID string, content string) (*types.Post, error) {
	path := fmt.Sprintf("/api/posts/%s/repost", postID)

	payload := repostRequest{
		Content: content,
	}

	req, err := s.transport.NewRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.Post

	err = json.UnmarshalRead(resp.Body, &result, transport.DataOptions)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Vote голосует в опросе, прикреплённом к посту.
func (s *Service) Vote(ctx context.Context, postID string, optionIDs ...string) (*types.Poll, error) {
	path := fmt.Sprintf("/api/posts/%s/poll/vote", postID)

	payload := voteRequest{
		OptionIds: optionIDs,
	}

	req, err := s.transport.NewRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.Poll `json:"data"`
	}

	err = json.UnmarshalRead(resp.Body, &result, transport.DataOptions)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// View отмечает пост как просмотренный.
func (s *Service) View(ctx context.Context, postID string) (*types.PostViewResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/view", postID)

	req, err := s.transport.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// API может вернуть 204 No Content для успешного просмотра
	if resp.StatusCode == 204 {
		return &types.PostViewResponse{Viewed: true}, nil
	}

	var result types.PostViewResponse
	err = json.UnmarshalRead(resp.Body, &result, transport.DataOptions)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
