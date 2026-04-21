package user

import "github.com/slipynil/itd-go/types"

type UpdateProfile struct {
	DisplayName string `json:"displayName,omitempty"`
	Username    string `json:"username,omitempty"`
	Bio         string `json:"bio,omitempty"`
	BannerID    string `json:"bannerId,omitempty"`
}

// ResponseUsers представляет ответ API при получении списка пользователей.
type ResponseUsers struct {
	// Data - данные со списком пользователей и информацией о пагинации
	Data UsersData `json:"data"`
}

// UsersData содержит массив пользователей и информацию о пагинации.
type UsersData struct {
	// Users - массив пользователей на текущей странице
	Users []types.UserCompact `json:"users"`

	// Pagination - информация о пагинации
	Pagination Pagination `json:"pagination"`
}

// Pagination содержит информацию о page-based пагинации для списков пользователей.
type Pagination struct {
	// Page - номер текущей страницы
	Page int `json:"page"`

	// Limit - количество пользователей на странице
	Limit int `json:"limit"`

	// Total - общее количество пользователей
	Total int `json:"total"`
}
