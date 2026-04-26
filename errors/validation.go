package errors

import "errors"

// Validation errors - клиентская валидация перед отправкой запроса

// ErrEmptyContent возвращается когда контент и файлы отсутствуют.
var ErrEmptyContent = errors.New("content or files required")

// ErrEmptyPostID возвращается при пустом ID поста.
var ErrEmptyPostID = errors.New("postID cannot be empty")

// ErrEmptyCommentID возвращается при пустом ID комментария.
var ErrEmptyCommentID = errors.New("commentID cannot be empty")

// ErrEmptyUserID возвращается при пустом ID пользователя.
var ErrEmptyUserID = errors.New("userID cannot be empty")

// ErrEmptyRefreshToken возвращается при пустом refresh token в конфиге.
var ErrEmptyRefreshToken = errors.New("refresh token is empty")

// Poll validation

// ErrNilPoll возвращается когда poll = nil при создании поста с опросом.
var ErrNilPoll = errors.New("poll cannot be nil")

// ErrInsufficientPollOptions возвращается когда в опросе < 2 вариантов.
var ErrInsufficientPollOptions = errors.New("poll must have at least 2 options")

// Reply validation

// ErrEmptyReplyToUserID возвращается при создании ответа без указания пользователя.
var ErrEmptyReplyToUserID = errors.New("replyToUserID cannot be empty")

// ErrEmptyParentCommentID возвращается при создании ответа без указания родительского комментария.
var ErrEmptyParentCommentID = errors.New("parentCommentID cannot be empty")
