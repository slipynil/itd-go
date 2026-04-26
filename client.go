package itdgo

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/slipynil/itd-go/api/comments"
	"github.com/slipynil/itd-go/api/posts"
	"github.com/slipynil/itd-go/api/user"
	"github.com/slipynil/itd-go/internal/root"
)

// ITD_DOMAIN - доменное имя ITD API в формате punycode.
const ITD_DOMAIN string = "xn--d1ah4a.com"

// BASE_URL - базовый URL для всех запросов к ITD API.
const BASE_URL string = "https://" + ITD_DOMAIN

// SDK_VERSION - текущая версия ITD Go SDK.
const SDK_VERSION string = "0.3.1"

// Client - главный клиент ITD SDK для взаимодействия с API.
// Предоставляет доступ к различным группам API методов.
type Client struct {
	// Posts - API для работы с постами
	Posts posts.Service

	// User - API для работы с пользователями
	User user.Service

	// Comments - API для работы с комментариями
	Comments comments.Service
}

// New создаёт и инициализирует новый экземпляр ITD клиента.
// Выполняет аутентификацию с использованием refresh token из конфигурации.
//
// Параметры:
//   - ctx: контекст для управления временем жизни операции инициализации
//   - cfg: конфигурация клиента (RefreshToken обязателен)
//
// Возвращает инициализированный клиент или ошибку при проблемах с аутентификацией.
func New(ctx context.Context, cfg Config) (*Client, error) {
	if !cfg.WithoutBanner {
		printBanner(SDK_VERSION, os.Stdout)
	}

	apiCfg := root.Config{
		RefreshToken: cfg.RefreshToken,
		Url:          BASE_URL,
		Domain:       ITD_DOMAIN,
		UserAgent:    cfg.UserAgent,
		Timeout:      cfg.Timeout,
	}
	root, err := root.New(ctx, apiCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create api client: %w", err)
	}

	return &Client{
		Posts:    *root.Posts,
		User:     *root.User,
		Comments: *root.Comments,
	}, nil
}

const (
	colorBlue  = "\033[34m"
	colorCyan  = "\033[36m"
	colorReset = "\033[0m"
	colorGray  = "\033[90m"
)

const bannerText = `
  ██╗████████╗██████╗      ██████╗  ██████╗
  ██║╚══██╔══╝██╔══██╗    ██╔════╝ ██╔═══██╗
  ██║   ██║   ██║  ██║    ██║  ███╗██║   ██║
  ██║   ██║   ██║  ██║    ██║   ██║██║   ██║
  ██║   ██║   ██████╔╝    ╚██████╔╝╚██████╔╝
  ╚═╝   ╚═╝   ╚═════╝      ╚═════╝  ╚═════╝
`

func printBanner(version string, w io.Writer) {
	fmt.Fprintf(w, "%s%s%s", colorCyan, bannerText, colorReset)
	fmt.Fprintf(w, "  %sITD Go SDK%s v%s %s— неофициальный клиент итд.com%s\n\n",
		colorBlue, colorReset, version, colorGray, colorReset)
}
