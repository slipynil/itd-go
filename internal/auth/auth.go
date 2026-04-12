package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/slipynil/itd-go/internal/dto"
	"github.com/slipynil/itd-go/internal/pkg/errors"
	pkg "github.com/slipynil/itd-go/internal/pkg/jwt"
)

// Client реализует Provider, используя refresh token для аутентификации.
// Автоматически обновляет access token при истечении срока действия.
type Client struct {
	mu          sync.RWMutex // мютекс для потокобезопасного доступа к токену
	BaseURL     string       // базовый URL для запросов аутентификации
	HttpClient  *http.Client // HTTP клиент для выполнения запросов
	accessToken string       // текущий access token
	tokenExpiry time.Time    // время истечения срока действия токена
	userID      string       // ID аутентифицированного пользователя
}

// New создаёт нового провайдера аутентификации через refresh token.
func New(cfg Config) (*Client, error) {

	auth := &Client{
		BaseURL:    cfg.Url,
		HttpClient: cfg.HttpClient,
	}

	return auth, nil
}

// GetToken возвращает access token, обновляя его при необходимости
func (r *Client) GetAccessToken(ctx context.Context) (string, error) {
	r.mu.RLock()
	token := r.accessToken
	expiry := r.tokenExpiry
	r.mu.RUnlock()

	// токен есть и ещё не истёк (с буфером 30 секунд)
	if token != "" && time.Now().Before(expiry.Add(-30*time.Second)) {
		return token, nil
	}

	// Нужно обновить токен
	return r.refreshAccessToken(ctx)
}

// GetUserID возвращает ID пользователя
func (r *Client) GetUserID() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.userID
}

// IsAuthenticated возвращает true если есть валидный токен
func (r *Client) IsAuthenticated() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.accessToken != "" && r.userID != ""
}

// возвращает jar cookies для клиента
func (c *Client) GetCookieJar() http.CookieJar {
	return c.HttpClient.Jar
}

// refreshAccessToken выполняет обновление access token
func (r *Client) refreshAccessToken(ctx context.Context) (string, error) {

	// Формируем URL для refresh endpoint
	endpoint, err := url.JoinPath("api", "v1", "auth", "refresh")
	if err != nil {
		return "", fmt.Errorf("ошибка при формировании endpoint: %w", err)
	}

	refreshURL, err := url.JoinPath(r.BaseURL, endpoint)
	if err != nil {
		return "", fmt.Errorf("ошибка при формировании URL: %w", err)
	}

	// Выполняем POST запрос
	req, err := http.NewRequestWithContext(ctx, "POST", refreshURL, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при обновлении токена: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if err := errors.CheckResponse(resp); err != nil {
		return "", err
	}

	// Декодируем ответ
	data := new(dto.AuthResponse)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return "", fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if data.AccessToken == "" {
		return "", fmt.Errorf("получен пустой access token")
	}

	expiry, err := pkg.ParseJWTExpiry(data.AccessToken)
	if err != nil {
		expiry = time.Now().Add(604800 * time.Second)
	}

	// Сохраняем токен и его время действия
	r.mu.Lock()
	r.tokenExpiry = expiry
	r.accessToken = data.AccessToken
	defer r.mu.Unlock()

	return r.accessToken, nil
}

// fetchUserID получает и сохраняет ID пользователя
func (r *Client) fetchUserID(ctx context.Context) error {
	// Формируем URL для /api/users/me
	endpoint, err := url.JoinPath("api", "users", "me")
	if err != nil {
		return fmt.Errorf("ошибка при формировании endpoint: %w", err)
	}

	meURL, err := url.JoinPath(r.BaseURL, endpoint)
	if err != nil {
		return fmt.Errorf("ошибка при формировании URL: %w", err)
	}

	// Создаем запрос с Authorization заголовком
	req, err := http.NewRequestWithContext(ctx, "GET", meURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.accessToken)

	// обрабатываем ответ
	resp, err := r.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при получении информации о пользователе: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if err := errors.CheckResponse(resp); err != nil {
		return err
	}

	// Декодируем ответ
	var userData dto.UserIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if userData.ID == "" {
		return fmt.Errorf("получен пустой user ID")
	}

	r.userID = userData.ID
	return nil
}
