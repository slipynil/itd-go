package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// httpTransport [Client] - это тупой слой, он знает только как отправить запрос и получить ответ.
// Он не знает что такое `User`, `Post`, `Feed`.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создает новый transport клиент с middleware цепочкой
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

func (c *Client) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("ошибка при декодировании тела запроса: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
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
