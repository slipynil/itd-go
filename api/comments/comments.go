package comments

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-json-experiment/json"

	"github.com/slipynil/itd-go/errors"
	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// Service предоставляет методы для работы с комментариями ITD API.
type Service struct {
	transport *transport.Client
}

// New создаёт новый экземпляр клиента для работы с комментариями.
func New(t *transport.Client) *Service {
	return &Service{transport: t}
}

// NewCommentList создаёт итератор для получения комментариев к посту.
// Параметры:
//   - postID: идентификатор поста
//   - limit: количество комментариев на страницу (рекомендуется 10-20)
//
// Возвращает CommentIterator для постраничной загрузки комментариев.
func (s *Service) NewCommentList(postID string, limit int) CommentIterator {
	return commentListIterator(s, postID, limit)
}

// ListReplies получает список ответов на комментарий.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - commentID: идентификатор комментария
//   - limit: количество ответов на страницу
//
// Возвращает массив ответов или ошибку при проблемах с сетью/API.
func (s *Service) ListReplies(ctx context.Context, commentID string, limit int) ([]*types.Comment, error) {
	result, err := s.getReplyList(ctx, commentID, limit)
	if err != nil {
		return nil, err
	}
	return result.Data.Replies, nil
}

func (s *Service) getReplyList(ctx context.Context, commentID string, limit int) (*repliesResponse, error) {
	path := fmt.Sprintf("/api/comments/%s/replies?limit=%d", commentID, limit)

	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result repliesResponse

	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateComment создаёт новый комментарий к посту.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - postID: идентификатор поста
//   - content: текстовое содержимое комментария
//   - filePaths: пути к файлам для загрузки и прикрепления к комментарию
//
// Возвращает созданный комментарий или ошибку при проблемах с сетью/API.
func (s *Service) CreateComment(ctx context.Context, postID string, content string, filePaths ...string) (*types.CreatedComment, error) {
	if postID == "" {
		return nil, errors.ErrEmptyPostID
	}
	if strings.TrimSpace(content) == "" && len(filePaths) == 0 {
		return nil, errors.ErrEmptyContent
	}
	if err := validateCommentFiles(filePaths); err != nil {
		return nil, err
	}

	attachmentIDs, err := s.transport.UploadFiles(ctx, filePaths...)
	if err != nil {
		return nil, err
	}

	payload := createCommentRequest{
		Content:       content,
		AttachmentIDs: attachmentIDs,
	}

	path := fmt.Sprintf("/api/posts/%s/comments", postID)
	req, err := s.transport.NewRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.CreatedComment
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateReply создаёт ответ на комментарий.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - parentCommentID: идентификатор родительского комментария
//   - replyToUserID: идентификатор пользователя, которому адресован ответ
//   - content: текстовое содержимое ответа
//   - filePaths: пути к файлам для загрузки и прикрепления к ответу
//
// Возвращает созданный ответ или ошибку при проблемах с сетью/API.
func (s *Service) CreateReply(
	ctx context.Context,
	parentCommentID,
	replyToUserID,
	content string,
	filePaths ...string,
) (*types.CreatedComment, error) {
	if parentCommentID == "" {
		return nil, errors.ErrEmptyParentCommentID
	}
	if replyToUserID == "" {
		return nil, errors.ErrEmptyReplyToUserID
	}
	if strings.TrimSpace(content) == "" && len(filePaths) == 0 {
		return nil, errors.ErrEmptyContent
	}
	if err := validateCommentFiles(filePaths); err != nil {
		return nil, err
	}

	attachmentIDs, err := s.transport.UploadFiles(ctx, filePaths...)
	if err != nil {
		return nil, err
	}

	payload := createReplyRequest{
		ReplyToUserId: replyToUserID,
		Content:       content,
		AttachmentIDs: attachmentIDs,
	}

	path := fmt.Sprintf("/api/comments/%s/replies", parentCommentID)
	req, err := s.transport.NewRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.CreatedComment
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete удаляет комментарий по его ID.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - commentID: идентификатор комментария для удаления
//
// Возвращает ошибку при проблемах с сетью/API или если комментарий не найден.
func (s *Service) Delete(ctx context.Context, commentID string) error {
	path := fmt.Sprintf("/api/comments/%s", commentID)

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

// Like ставит лайк на комментарий.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - commentID: идентификатор комментария
//
// Возвращает ошибку при проблемах с сетью/API.
func (s *Service) Like(ctx context.Context, commentID string) error {
	path := fmt.Sprintf("/api/comments/%s/like", commentID)

	req, err := s.transport.NewRequest(ctx, "POST", path, nil)
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

// Unlike убирает лайк с комментария.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - commentID: идентификатор комментария
//
// Возвращает ошибку при проблемах с сетью/API.
func (s *Service) Unlike(ctx context.Context, commentID string) error {
	path := fmt.Sprintf("/api/comments/%s/like", commentID)

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

// Update обновляет содержимое комментария.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - commentID: идентификатор комментария
//   - content: новое текстовое содержимое
//
// Возвращает обновлённый комментарий или ошибку при проблемах с сетью/API.
func (s *Service) Update(ctx context.Context, commentID string, content string) (*types.CommentUpdate, error) {
	path := fmt.Sprintf("/api/comments/%s", commentID)

	payload := updateCommentRequest{
		Content: content,
	}

	req, err := s.transport.NewRequest(ctx, "PATCH", path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
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

func validateCommentFiles(paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	if len(paths) > 4 {
		return fmt.Errorf("%w: %d, max: 10", errors.TooManyFiles, len(paths))
	}

	allowed := []string{".png", ".webp", ".ogg", ".mp3"}
	for _, path := range paths {
		ext := strings.ToLower(filepath.Ext(path))
		if ext == "" {
			return fmt.Errorf("%w: %s", errors.NoFileExtension, path)
		}
		if !slices.Contains(allowed, ext) {
			return fmt.Errorf("%w: %s, supported: %v", errors.InvalidFileExtension, ext, allowed)
		}
	}
	return nil
}
