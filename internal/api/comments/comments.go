package comments

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// Comments предоставляет методы для работы с комментариями ITD API.
type Comments struct {
	transport *transport.Client
}

// New создаёт новый экземпляр клиента для работы с комментариями.
func New(t *transport.Client) *Comments {
	return &Comments{transport: t}
}

func (c *Comments) NewCommentList(ctx context.Context, postID string, limit int) types.CommentIterator {
	return commentListIterator(c, ctx, postID, limit)
}

// getCommentList получает сырую json стурктуру с комментариями и информацией о пагинации.
func (c *Comments) getCommentList(ctx context.Context, postID, cursor string, limit int) (*CommentsResponse, error) {
	path := fmt.Sprintf("/api/posts/%s/comments?limit=%d&sort=popular", postID, limit)

	req, err := c.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CommentsResponse
	if err := json.UnmarshalRead(resp.Body, &result, transport.DataOptions); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListReplies получает список ответов на комментарий.
func (c *Comments) ListReplies(ctx context.Context, commentID string, limit int) ([]*types.Comment, error) {
	result, err := c.getReplyList(ctx, commentID, limit)
	if err != nil {
		return nil, err
	}
	return result.Data.Replies, nil
}

func (c *Comments) getReplyList(ctx context.Context, commentID string, limit int) (*RepliesResponse, error) {
	path := fmt.Sprintf("/api/comments/%s/replies?limit=%d", commentID, limit)

	req, err := c.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RepliesResponse

	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create создаёт новый комментарий к посту.
func (c *Comments) CreateComment(ctx context.Context, postID string, content string, attachmentIDs []string) (*types.CreateComment, error) {
	path := fmt.Sprintf("/api/posts/%s/comments", postID)

	if attachmentIDs == nil {
		attachmentIDs = []string{}
	}

	payload := map[string]interface{}{
		"content":       content,
		"attachmentIds": attachmentIDs,
	}

	req, err := c.transport.NewRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.CreateComment
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Reply создаёт ответ на комментарий.
func (c *Comments) CreateReply(
	ctx context.Context,
	parentCommentID,
	replyToUserID,
	content string,
	attachmentIDs []string,
) (*types.CreateComment, error) {
	path := fmt.Sprintf("/api/comments/%s/replies", parentCommentID)

	if attachmentIDs == nil {
		attachmentIDs = []string{}
	}

	payload := map[string]interface{}{
		"replyToUserId": replyToUserID,
		"content":       content,
		"attachmentIds": attachmentIDs,
	}

	req, err := c.transport.NewRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.CreateComment
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete удаляет комментарий по его ID.
func (c *Comments) Delete(ctx context.Context, commentID string) error {
	path := fmt.Sprintf("/api/comments/%s", commentID)

	req, err := c.transport.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Like ставит лайк на комментарий.
func (c *Comments) Like(ctx context.Context, commentID string) error {
	path := fmt.Sprintf("/api/comments/%s/like", commentID)

	req, err := c.transport.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Unlike убирает лайк с комментария.
func (c *Comments) Unlike(ctx context.Context, commentID string) error {
	path := fmt.Sprintf("/api/comments/%s/like", commentID)

	req, err := c.transport.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Update обновляет содержимое комментария.
func (c *Comments) Update(ctx context.Context, commentID string, content string) (*types.CommentUpdate, error) {
	path := fmt.Sprintf("/api/comments/%s", commentID)

	payload := map[string]interface{}{
		"content": content,
	}

	req, err := c.transport.NewRequest(ctx, "PATCH", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.CommentUpdate
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
