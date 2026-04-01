package root

import (
	"context"
	"fmt"

	"github.com/slipynil/itd-go/internal/api/posts"
	"github.com/slipynil/itd-go/internal/api/user"
	"github.com/slipynil/itd-go/internal/auth"
	"github.com/slipynil/itd-go/internal/transport"
)

type Client struct {
	Posts *posts.Posts
	User  *user.User
}

// New создает новый [Client] с дефолтной конфигурацией.
// Поле refreshToken обязательное, остальные могут быть пустыми.
//
// Перед использованием убедитесь, что refreshToken валиден.
func New(ctx context.Context, cfg Config) (*Client, error) {
	// Проверка refresh token
	if cfg.RefreshToken == "" {
		return nil, fmt.Errorf("поле refreshToken пустое")
	}

	httpClient, err := CreateHttpClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка в создании http client: %w", err)
	}

	// Создаем auth провайдер
	authCfg := auth.Config{
		Url:        cfg.Url,
		HttpClient: httpClient,
	}
	authClient, err := auth.New(authCfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка в создании auth provider: %w", err)
	}

	// Выполняем начальную аутентификацию
	if _, err := authClient.GetAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("ошибка в аутентификации: %w", err)
	}

	// Создаем transport клиент
	transportCfg := transport.Config{
		BaseURL:    cfg.Url,
		HttpClient: httpClient,
		AuthClient: authClient,
	}
	transportClient := transport.NewClient(transportCfg)

	posts := posts.New(transportClient)
	user := user.New(transportClient)
	return &Client{
		Posts: posts,
		User:  user,
	}, nil
}
