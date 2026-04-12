package types

// Pin представляет пин — награду или достижение пользователя.
type Pin struct {
	// Slug - уникальный идентификатор пина
	Slug string `json:"slug"`

	// Name - отображаемое название пина
	Name string `json:"name"`

	// Description - описание за что выдан пин
	Description string `json:"description"`

	// URL - ссылка на анимацию или изображение пина
	URL string `json:"url"`
}
