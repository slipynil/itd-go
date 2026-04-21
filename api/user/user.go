package user

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// User предоставляет методы для работы с пользователями ITD API.
type Service struct {
	transport *transport.Client
}

// New создаёт новый экземпляр клиента для работы с пользователями.
func New(t *transport.Client) *Service {
	return &Service{transport: t}
}

// Me получает информацию о текущем аутентифицированном пользователе.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//
// Возвращает информацию о текущем пользователе или ошибку при проблемах с аутентификацией.
func (s *Service) Me(ctx context.Context) (*types.Me, error) {
	req, err := s.transport.NewRequest(ctx, "GET", "/api/users/me", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.Me

	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get получает информацию о пользователе по его username.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - username: имя пользователя (без @)
//
// Возвращает полную информацию о пользователе или ошибку при проблемах с сетью/API.
func (s *Service) Get(ctx context.Context, username string) (*types.User, error) {
	path := fmt.Sprintf("/api/users/%s", username)
	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.User

	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Follow подписывается на пользователя.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - username: имя пользователя для подписки
//
// Возвращает ошибку при проблемах с сетью/API или если пользователь не найден.
func (s *Service) Follow(ctx context.Context, username string) error {
	path := fmt.Sprintf("/api/users/%s/follow", username)
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

// Unfollow отписывается от пользователя.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - username: имя пользователя для отписки
//
// Возвращает ошибку при проблемах с сетью/API или если пользователь не найден.
func (s *Service) Unfollow(ctx context.Context, username string) error {
	path := fmt.Sprintf("/api/users/%s/follow", username)
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

// UpdateProfile обновляет профиль текущего пользователя.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - config: структура с полями для обновления (пустые поля не изменяются)
//
// Возвращает обновлённую информацию о пользователе или ошибку при проблемах с сетью/API.
func (s *Service) UpdateProfile(ctx context.Context, config types.UpdateProfile) (*types.UpdateProfileResponse, error) {

	payload := UpdateProfile{
		DisplayName: config.DisplayName,
		Username:    config.Username,
		Bio:         config.Bio,
	}

	if config.BannerPath != "" {
		banner, err := s.transport.Upload(ctx, config.BannerPath)
		if err != nil {
			return nil, err
		}
		payload.BannerID = banner.ID
	}

	req, err := s.transport.NewRequest(ctx, "PUT", "/api/users/me", payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.UpdateProfileResponse

	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetFollowers возвращает список подписчиков пользователя.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - username: идентификатор пользователя или юзернейм пользователя
//   - limit: максимальное количество подписчиков в ответе
//
// Возвращает срез UserCompact с данными подписчиков.
//
// Примечание: в настоящее время API возвращает не более 20 подписчиков
// независимо от переданного limit.
func (s *Service) GetFollowers(ctx context.Context, username string, limit int) ([]types.UserCompact, error) {
	result, err := s.getFollowers(ctx, username, limit, 1)
	if err != nil {
		return nil, err
	}
	return result.Users, nil
}

// GetFollowing возвращает список подписок пользователя.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - username: идентификатор пользователя или юзернейм пользователя
//   - limit: максимальное количество подписок в ответе
//
// Возвращает срез UserCompact с данными подписок.
//
// Примечание: в настоящее время API возвращает не более 20 подписок
// независимо от переданного limit.
func (s *Service) GetFollowing(ctx context.Context, username string, limit int) ([]types.UserCompact, error) {
	result, err := s.getFollowing(ctx, username, limit, 1)
	if err != nil {
		return nil, err
	}
	return result.Users, nil
}

// getFollowers получает подписчиков с page-based пагинацией (внутренний метод для итератора).
func (s *Service) getFollowers(ctx context.Context, username string, limit, page int) (*UsersData, error) {

	path := fmt.Sprintf("/api/users/%s/followers?limit=%d&page=%d", username, limit, page)

	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ResponseUsers
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// getFollowing получает подписки с page-based пагинацией (внутренний метод для итератора).
func (s *Service) getFollowing(ctx context.Context, userID string, limit, page int) (*UsersData, error) {

	path := fmt.Sprintf("/api/users/%s/following?limit=%d&page=%d", userID, limit, page)

	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ResponseUsers
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
