package root

import (
	"context"
	"fmt"

	"github.com/slipynil/itd-go/api/comments"
	"github.com/slipynil/itd-go/api/notifications"
	"github.com/slipynil/itd-go/api/posts"
	"github.com/slipynil/itd-go/api/user"
	"github.com/slipynil/itd-go/errors"
	"github.com/slipynil/itd-go/internal/auth"
	"github.com/slipynil/itd-go/internal/transport"
)

// Client предоставляет доступ ко всем API модулям ITD SDK.
// Содержит клиенты для работы с постами, пользователями и комментариями.
type Client struct {
	// Posts - клиент для работы с постами
	Posts *posts.Service

	// User - клиент для работы с пользователями
	User *user.Service

	// Comments - клиент для работы с комментариями
	Comments *comments.Service

	// Notifications - клиент для работы с уведомлениями
	Notifications *notifications.Service
}

// New создаёт новый корневой клиент SDK с настроенной аутентификацией.
// RefreshToken в конфигурации является обязательным параметром.
func New(ctx context.Context, cfg Config) (*Client, error) {
	// Проверка refresh token
	if cfg.RefreshToken == "" {
		return nil, errors.ErrEmptyRefreshToken
	}

	httpClient, err := CreateHttpClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	// Создаем auth провайдер
	authCfg := auth.Config{
		Url:        cfg.Url,
		HttpClient: httpClient,
	}
	authClient, err := auth.New(authCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth provider: %w", err)
	}

	// Выполняем начальную аутентификацию
	if _, err := authClient.GetAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Создаем transport клиент
	transportCfg := transport.Config{
		BaseURL:    cfg.Url,
		HttpClient: httpClient,
		AuthClient: authClient,
	}
	t := transport.NewClient(transportCfg)

	posts := posts.New(t)
	user := user.New(t)
	comments := comments.New(t)
	notifications := notifications.New(t)

	return &Client{
		Posts:         posts,
		User:          user,
		Comments:      comments,
		Notifications: notifications,
	}, nil
}
