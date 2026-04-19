package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client предоставляет низкоуровневый HTTP транспорт для выполнения запросов к API.
// Не содержит бизнес-логики, только отправку запросов и получение ответов.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создаёт новый transport клиент с настроенной middleware цепочкой.
func NewClient(cfg Config) *Client {
	transport := buildTransport(cfg)

	return &Client{
		httpClient: &http.Client{
			Transport: transport,
			Jar:       cfg.AuthClient.GetCookieJar(),
		},
		baseURL: cfg.BaseURL,
	}
}

// Do выполняет HTTP-запрос с контекстом
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// NewRequestMultipart создаёт HTTP-запрос с multipart/form-data телом.
// Используется для загрузки файлов на сервер.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - method: HTTP метод (обычно POST)
//   - path: путь к API endpoint
//   - body: буфер с multipart данными
//   - contentType: Content-Type заголовок (обычно из writer.FormDataContentType())
// Возвращает готовый HTTP-запрос или ошибку.
func (c *Client) NewRequestMultipart(
	ctx context.Context,
	method,
	path string,
	body *bytes.Buffer,
	contentType string,
) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	return req, nil
}

func (c *Client) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации тела запроса: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// buildTransport создает цепочку middleware
func buildTransport(cfg Config) http.RoundTripper {
	// Используем Transport из httpClient как базовый
	base := cfg.HttpClient.Transport
	if base == nil {
		base = http.DefaultTransport
	}

	// Цепочка: base -> statusCheck -> auth
	var transport http.RoundTripper = &statusCheckMiddleware{
		base: base,
	}

	if cfg.AuthClient != nil {
		transport = &authMiddleware{
			base:     transport,
			provider: cfg.AuthClient,
		}
	}

	return transport
}

// обработка ответов
func (c *Client) DoJSON(ctx context.Context, req *http.Request, result any) error {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверка статуса теперь выполняется в statusCheckMiddleware

	if result == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}
	return nil
}
