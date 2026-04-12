package user

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
)

// User предоставляет методы для работы с пользователями ITD API.
type User struct {
	transport *transport.Client
}

// New создаёт новый экземпляр клиента для работы с пользователями.
func New(t *transport.Client) *User {
	return &User{transport: t}
}

// Me получает информацию о текущем аутентифицированном пользователе.
func (u *User) Me(ctx context.Context) (*types.Me, error) {
	req, err := u.transport.NewRequest(ctx, "GET", "/api/users/me", nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.transport.Do(req)
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
func (u *User) Get(ctx context.Context, username string) (*types.User, error) {
	path := fmt.Sprintf("/api/users/%s", username)
	req, err := u.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.transport.Do(req)
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
func (u *User) Follow(ctx context.Context, username string) error {
	path := fmt.Sprintf("/api/users/%s/follow", username)
	req, err := u.transport.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return err
	}

	resp, err := u.transport.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Unfollow отписывается от пользователя.
func (u *User) Unfollow(ctx context.Context, username string) error {
	path := fmt.Sprintf("/api/users/%s/follow", username)
	req, err := u.transport.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := u.transport.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UpdateProfile обновляет профиль текущего пользователя.
func (u *User) UpdateProfile(ctx context.Context, displayName *string, username *string, bio *string, bannerID *string) (*types.UpdateProfileResponse, error) {
	payload := make(map[string]string)

	if displayName != nil {
		payload["displayName"] = *displayName
	}
	if username != nil {
		payload["username"] = *username
	}
	if bio != nil {
		payload["bio"] = *bio
	}
	if bannerID != nil {
		fmt.Println("в этой версии SDK баннер пока что нельзя менять")
		payload["bannerId"] = *bannerID
	}

	req, err := u.transport.NewRequest(ctx, "PUT", "/api/users/me", payload)
	if err != nil {
		return nil, err
	}

	resp, err := u.transport.Do(req)
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

func (u *User) GetFollowers(ctx context.Context, username string, limit int) ([]types.UserCompact, error) {
	result, err := u.getFollowers(ctx, username, limit, 1)
	if err != nil {
		return nil, err
	}
	return result.Users, nil
}

func (u *User) GetFollowing(ctx context.Context, username string, limit int) ([]types.UserCompact, error) {
	result, err := u.getFollowing(ctx, username, limit, 1)
	if err != nil {
		return nil, err
	}
	return result.Users, nil
}

// getFollowers получает подписчиков с page-based пагинацией (внутренний метод для итератора).
func (u *User) getFollowers(ctx context.Context, username string, limit, page int) (*UsersData, error) {

	path := fmt.Sprintf("/api/users/%s/followers?limit=%d&page=%d", username, limit, page)

	req, err := u.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.transport.Do(req)
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
func (u *User) getFollowing(ctx context.Context, userID string, limit, page int) (*UsersData, error) {

	path := fmt.Sprintf("/api/users/%s/following?limit=%d&page=%d", userID, limit, page)

	req, err := u.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.transport.Do(req)
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
