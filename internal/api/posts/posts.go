package posts

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// Posts предоставляет методы для работы с постами ITD API.
type Posts struct {
	transport *transport.Client
}

// New создаёт новый экземпляр клиента для работы с постами.
func New(t *transport.Client) *Posts {
	return &Posts{transport: t}
}

// NewFeed возвращает итератор для получения постов.
func (p *Posts) NewFeed(ctx context.Context, tab types.FeedTab, limit int) types.FeedIterator {
	return newFeedIterator(p, ctx, tab, limit)
}

// NewUserPosts возвращает итератор для получения постов пользователя.
func (p *Posts) NewUserPosts(ctx context.Context, username string, limit int) types.FeedIterator {
	return newUserPostsIterator(p, ctx, username, limit)
}

// Get получает пост по его ID.
func (p *Posts) Get(ctx context.Context, postID string) (*types.Post, error) {
	path := fmt.Sprintf("/api/posts/%s", postID)
	req, err := p.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.transport.Do(req)
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
func (p *Posts) Create(ctx context.Context, content string, filePaths ...string) (*types.Post, error) {

	attachmentIDs := make([]string, 0, len(filePaths))

	for _, path := range filePaths {
		attachment, err := p.transport.Upload(ctx, path)
		if err != nil {
			return nil, err
		}
		attachmentIDs = append(attachmentIDs, attachment.ID)
	}

	payload := createPostRequest{
		Content:       content,
		AttachmentIDs: attachmentIDs,
	}

	req, err := p.transport.NewRequest(ctx, "POST", "/api/posts", payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.transport.Do(req)
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
func (p *Posts) CreateWithPoll(ctx context.Context, content string, poll *types.PollRequest, filePaths ...string) (*types.Post, error) {

	attachmentIDs := make([]string, 0, len(filePaths))

	for _, path := range filePaths {
		attachment, err := p.transport.Upload(ctx, path)
		if err != nil {
			return nil, err
		}
		attachmentIDs = append(attachmentIDs, attachment.ID)
	}

	payload := createPostRequest{
		Content:       content,
		AttachmentIDs: attachmentIDs,
		Poll:          poll,
	}

	req, err := p.transport.NewRequest(ctx, "POST", "/api/posts", payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.transport.Do(req)
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
func (p *Posts) Delete(ctx context.Context, postID string) error {
	path := fmt.Sprintf("/api/posts/%s", postID)
	req, err := p.transport.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := p.transport.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Like ставит лайк на пост.
func (p *Posts) Like(ctx context.Context, postID string) (*types.LikesCountResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/like", postID)
	req, err := p.transport.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.transport.Do(req)
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
func (p *Posts) Unlike(ctx context.Context, postID string) (*types.LikesCountResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/like", postID)
	req, err := p.transport.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.transport.Do(req)
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
func (p *Posts) Repost(ctx context.Context, postID string, content string) (*types.Post, error) {
	path := fmt.Sprintf("/api/posts/%s/repost", postID)

	var payload map[string]interface{}
	if len(content) != 0 {
		payload = map[string]interface{}{
			"content": content,
		}
	}

	req, err := p.transport.NewRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.transport.Do(req)
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
func (p *Posts) Vote(ctx context.Context, postID string, optionIDs ...string) (*types.Poll, error) {
	path := fmt.Sprintf("/api/posts/%s/poll/vote", postID)

	payload := map[string]interface{}{}
	if len(optionIDs) > 0 {
		payload["optionIds"] = optionIDs
	}

	req, err := p.transport.NewRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.transport.Do(req)
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
func (p *Posts) View(ctx context.Context, postID string) (*types.PostViewResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/view", postID)

	req, err := p.transport.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.transport.Do(req)
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
