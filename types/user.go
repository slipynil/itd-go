package types

import "time"

// Me представляет полную информацию о текущем пользователе.
// Содержит профильные данные, статистику и настройки приватности.
type Me struct {
	// ID - уникальный идентификатор пользователя
	ID string `json:"id"`

	// Username - уникальный логин пользователя
	Username string `json:"username"`

	// DisplayName - отображаемое имя пользователя
	DisplayName string `json:"displayName"`

	// Avatar - аватар пользователя: эмодзи или URL изображения
	Avatar string `json:"avatar"`

	// BannerURL - URL баннера профиля (nil, если не установлен)
	BannerURL *string `json:"banner"`

	// Bio - текст биографии профиля (nil, если не заполнена)
	Bio *string `json:"bio"`

	// IsVerified - true, если аккаунт верифицирован
	IsVerified bool `json:"verified"`

	// IsPhoneVerified - true, если номер телефона подтверждён
	IsPhoneVerified bool `json:"isPhoneVerified"`

	// Pin - закреплённый пост (nil, если не установлен)
	Pin *Pin `json:"pin"`

	// IsPrivate - true, если профиль приватный
	IsPrivate bool `json:"isPrivate"`

	// WallAccess - кто может писать на стену ("everyone", "followers" и т.д.)
	WallAccess string `json:"wallAccess"`

	// LikesVisibility - кто может видеть лайки ("everyone", "followers" и т.д.)
	LikesVisibility string `json:"likesVisibility"`

	// IsFollowing - true, если текущий пользователь подписан на этот аккаунт
	IsFollowing bool `json:"isFollowing"`

	// IsFollowedBy - true, если этот аккаунт подписан на текущего пользователя
	IsFollowedBy bool `json:"isFollowedBy"`

	// FollowersCount - количество подписчиков
	FollowersCount int `json:"followersCount"`

	// FollowingCount - количество подписок
	FollowingCount int `json:"followingCount"`

	// PostsCount - количество постов пользователя
	PostsCount int `json:"postsCount"`

	// CreatedAt - дата регистрации аккаунта
	CreatedAt time.Time `json:"createdAt"`

	// Subscription - информация о подписке пользователя
	Subscription *Subscription `json:"subscription"`
}

// User представляет профиль пользователя.
type User struct {
	// ID - уникальный идентификатор пользователя
	ID string `json:"id"`

	// Username - уникальный логин пользователя
	Username string `json:"username"`

	// DisplayName - отображаемое имя пользователя
	DisplayName string `json:"displayName"`

	// Avatar - аватар пользователя: эмодзи или URL изображения
	Avatar string `json:"avatar"`

	// BannerURL - URL баннера профиля (nil, если не установлен)
	BannerURL *string `json:"banner"`

	// Bio - текст биографии профиля (nil, если не заполнена)
	Bio *string `json:"bio"`

	// IsVerified - true, если аккаунт верифицирован
	IsVerified bool `json:"verified"`

	// IsPhoneVerified - true, если номер телефона подтверждён (только для своего профиля)
	IsPhoneVerified bool `json:"isPhoneVerified"`

	// Pin - закреплённый пин пользователя (nil, если не установлен)
	Pin *Pin `json:"pin"`

	// HasNuksta - true, если пользователь поддержал ИТД (имеет подписку Нукста)
	HasNuksta bool `json:"hasNuksta"`

	// PinnedPostID - ID закреплённого поста (nil, если нет закреплённого поста)
	PinnedPostID *string `json:"pinnedPostId"`

	// IsPrivate - true, если профиль приватный
	IsPrivate bool `json:"isPrivate"`

	// WallAccess - кто может писать на стену ("everyone", "followers", "nobody" и т.д.)
	WallAccess string `json:"wallAccess"`

	// LikesVisibility - кто может видеть лайки ("everyone", "followers", "nobody" и т.д.)
	LikesVisibility string `json:"likesVisibility"`

	// IsFollowing - true, если текущий пользователь подписан на этот аккаунт
	IsFollowing bool `json:"isFollowing"`

	// IsFollowedBy - true, если этот аккаунт подписан на текущего пользователя
	IsFollowedBy bool `json:"isFollowedBy"`

	// FollowersCount - количество подписчиков
	FollowersCount int `json:"followersCount"`

	// FollowingCount - количество подписок
	FollowingCount int `json:"followingCount"`

	// PostsCount - количество постов пользователя
	PostsCount int `json:"postsCount"`

	// CreatedAt - дата регистрации аккаунта
	CreatedAt time.Time `json:"createdAt"`

	// Online - true, если пользователь сейчас онлайн
	Online bool `json:"online"`

	// LastSeen - информация о последнем визите пользователя (nil, если скрыто)
	LastSeen *LastSeen `json:"lastSeen"`
}

// UserCompact представляет краткую информацию о пользователе в списках
// (подписчики, подписки, поиск и т.д.)
type UserCompact struct {
	// ID - уникальный идентификатор пользователя
	ID string `json:"id"`

	// Username - уникальный логин пользователя
	Username string `json:"username"`

	// DisplayName - отображаемое имя пользователя
	DisplayName string `json:"displayName"`

	// Avatar - аватар пользователя: эмодзи или URL изображения
	Avatar string `json:"avatar"`

	// IsVerified - true, если аккаунт верифицирован
	IsVerified bool `json:"verified"`

	// IsFollowing - true, если текущий пользователь подписан на этот аккаунт
	IsFollowing bool `json:"isFollowing"`
}

// LastSeen представляет информацию о последнем визите пользователя.
type LastSeen struct {
	// Unit - относительное время последнего визита ("recently", "today", "week" и т.д.)
	Unit string `json:"unit"`
}

// Subscription представляет статус платной подписки пользователя.
type Subscription struct {
	// IsActive - true, если подписка активна
	IsActive bool `json:"isActive"`

	// ExpiresAt - дата окончания подписки (nil, если не активна)
	ExpiresAt *time.Time `json:"expiresAt"`

	// AutoRenewal - true, если включено автопродление подписки
	AutoRenewal bool `json:"autoRenewal"`
}

// UpdateProfileResponse представляет ответ API после обновления профиля.
// Содержит только изменяемые поля, а не полный профиль.
type UpdateProfileResponse struct {
	// ID - уникальный идентификатор пользователя
	ID string `json:"id"`

	// Username - уникальный логин пользователя
	Username string `json:"username"`

	// DisplayName - отображаемое имя пользователя
	DisplayName string `json:"displayName"`

	// Bio - текст биографии профиля (nil, если не заполнена)
	Bio *string `json:"bio"`

	// UpdatedAt - дата и время последнего обновления профиля
	UpdatedAt time.Time `json:"updatedAt"`
}
