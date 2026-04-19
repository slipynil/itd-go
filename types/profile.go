package types

// UpdateProfile представляет параметры для обновления профиля пользователя.
// Пустые поля не будут изменены на сервере.
type UpdateProfile struct {
	// DisplayName - новое отображаемое имя (пустая строка = не изменять)
	DisplayName string

	// Username - новое имя пользователя (пустая строка = не изменять)
	Username string

	// Bio - новая биография (пустая строка = не изменять)
	Bio string

	// BannerPath - путь к файлу баннера для автоматической загрузки (пустая строка = не изменять)
	BannerPath string
}
