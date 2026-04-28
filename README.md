<div align="center">

<img width="400" alt="Image" src="https://github.com/user-attachments/assets/8517c4b6-fd5f-498d-befc-8896976453dd" />


**Неофициальный Go SDK для работы с API социальной сети** [итд.com](https://итд.com)

**Версия:** 0.4.0

Made with ❤️ by [@Slipynil](https://github.com/slipynil)

[![Go Reference](https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/slipynil/itd-go)
[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/slipynil_chan)

</div>

# itd-go: Go SDK для API социальной сети итд.com

Package itd-go — библиотека для работы с API социальной сети [итд.com](https://итд.com).

Документация библиотеки доступна на [pkg.go.dev](https://pkg.go.dev/github.com/slipynil/itd-go).

> [!IMPORTANT]
> Это неофициальный SDK. API итд.com не документировано публично и может измениться без предупреждения.
> Используйте на свой риск. Авторы не несут ответственности за возможные блокировки аккаунта или другие последствия использования.

## Installation

```sh
go get github.com/slipynil/itd-go
```

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	itd "github.com/slipynil/itd-go"
	"github.com/slipynil/itd-go/types"
)

func main() {
	ctx := context.Background()

	client, err := itd.New(ctx, itd.Config{
		RefreshToken: "your_refresh_token",
		UserAgent:    "Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/149.0",
	})
	if err != nil {
		log.Fatalf("error creating client: %s\n", err)
	}

	// Create a post with text formatting and automatic file upload
	builder := types.NewPost("Hello from itd-go!").Bold("itd-go")
	post, err := client.Posts.Create(ctx, builder, "/path/to/image.jpg")
	if err != nil {
		log.Fatalf("error creating post: %s\n", err)
	}
	fmt.Printf("Post created: %s\n", post.ID)

	// Iterate through feed
	feed := client.Posts.NewFeed(types.FeedTabPopular, 20)
	for feed.HasMore() {
		posts, err := feed.Next(ctx)
		if err != nil {
			log.Fatalf("error fetching feed: %s\n", err)
		}
		for _, p := range posts {
			fmt.Printf("%s: %s\n", p.Author.DisplayName, p.Content)
		}
	}
}
```

## Configuration

### Getting Refresh Token

1. Откройте DevTools в браузере (F12)
2. Перейдите на вкладку Application/Storage → Cookies
3. Найдите cookie с именем `refresh_token` для домена `итд.com`
4. Скопируйте значение

### Config Options

```go
type Config struct {
	RefreshToken  string        // обязательно: refresh token для аутентификации
	UserAgent     string        // опционально: User-Agent для запросов
	Timeout       time.Duration // опционально: таймаут HTTP запросов (по умолчанию 30s)
	WithoutBanner bool          // опционально: отключить баннер при инициализации
}
```

## Features

- **Posts API**: создание, удаление, лайки, репосты, голосование в опросах, пагинация ленты
- **Text Formatting**: форматирование текста постов через PostBuilder (жирный, курсив, подчёркивание, зачёркивание, спойлер, моноширинный, ссылки)
- **User API**: профили пользователей, подписки, обновление профиля
- **Comments API**: комментарии и ответы, лайки, редактирование, пагинация
- **Notifications API**: получение уведомлений, пометка как прочитанных, real-time стрим через SSE
- **Automatic File Upload**: автоматическая загрузка файлов при создании постов и комментариев
- **Iterator Pattern**: удобная пагинация через итераторы для всех списочных методов
- **Token Management**: автоматическое обновление access token из refresh token

Полный список методов доступен в [документации](https://pkg.go.dev/github.com/slipynil/itd-go).

## Go Version Support

Библиотека требует Go 1.26 или выше.

## Examples

Примеры использования всех API методов находятся в директории [`examples/`](./examples/):
- [`examples/posts/`](./examples/posts/) — работа с постами
- [`examples/user/`](./examples/user/) — работа с пользователями
- [`examples/comments/`](./examples/comments/) — работа с комментариями
- [`examples/notifications/`](./examples/notifications/) — работа с уведомлениями

## License

[MIT](LICENSE)
