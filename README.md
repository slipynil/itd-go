# ITD Go SDK

Неофициальный Go SDK для работы с API социальной сети [итд.com](https://итд.com)

## Установка

```bash
go get github.com/slipynil/itd-go
```

## Быстрый старт

```go
package main

import (
    "context"
    "log"

    itd "github.com/slipynil/itd-go"
)

func main() {
    cfg := itd.Config{
        RefreshToken: "ваш_refresh_token",
        UserAgent:    "Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/149.0",
    }

    ctx := context.Background()
    client, err := itd.New(ctx, cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Создать итератор для ленты постов
    feed := client.Posts.NewFeed(ctx, 10, "popular")

    // Получить первую страницу постов
    posts, err := feed.Next()
    if err != nil {
        log.Fatal(err)
    }

    // Проверить, есть ли ещё посты
    if feed.HasMore() {
        // Получить следующую страницу
        morePosts, err := feed.Next()
        if err != nil {
            log.Fatal(err)
        }
        _ = morePosts
    }
}
```

## Конфигурация

```go
type Config struct {
    RefreshToken  string        // обязательно: refresh token для аутентификации
    UserAgent     string        // опционально: User-Agent для запросов
    Timeout       time.Duration // опционально: таймаут для HTTP запросов (по умолчанию 30s)
    WithoutBanner bool          // опционально: отключить вывод баннера при инициализации (по умолчанию false)
}
```

### Получение Refresh Token

1. Откройте DevTools в браузере (F12)
2. Перейдите на вкладку Application/Storage → Cookies
3. Найдите cookie с именем `refresh_token` для домена `итд.com`
4. Скопируйте значение

## Доступные методы

### Posts

Работа с постами через паттерн Iterator:

```go
// Создать итератор для ленты постов
// limit - количество постов на страницу (рекомендуется 10-50)
// sort - тип сортировки: "popular", "new", "following"
feed := client.Posts.NewFeed(ctx, limit, sort)

// Проверить, есть ли ещё посты для загрузки
hasMore := feed.HasMore()

// Получить следующую страницу постов
posts, err := feed.Next()
```

**Пример использования итератора:**

```go
feed := client.Posts.NewFeed(ctx, 20, "popular")

for feed.HasMore() {
    posts, err := feed.Next()
    if err != nil {
        log.Fatal(err)
    }

    for _, post := range posts {
        fmt.Printf("%s: %s\n", post.Author.DisplayName, post.Content)
    }
}
```

### User

API для работы с пользователями находится в разработке.

## Архитектура

SDK построен на многоуровневой архитектуре с чёткими границами ответственности:

```
itd-go/
├── client.go           — публичная точка входа
├── config.go           — публичная конфигурация
│
└── internal/
    ├── root/           — создание HTTP клиента, установка заголовков
    ├── auth/           — управление токенами, аутентификация
    ├── transport/      — HTTP механика, middleware
    └── api/            — группы методов (posts, user, и т.д.)
```

### Поток запроса

```
client.Posts.NewFeed(ctx, limit, sort)
    → создание feedIterator
    → iterator.Next()
        → transport.NewRequest()
        → transport.Do()
            → HTTP клиент с middleware цепочкой:
                1. Root Transport (Origin, Referer, User-Agent)
                2. Status Check Middleware (проверка статуса ≥400)
                3. Auth Middleware (добавление Authorization токена)
            → отправка запроса
            → декодирование JSON ответа
        → обновление cursor и hasMore
        → возврат постов
```

### Middleware цепочка

1. **Root Transport** — устанавливает базовые заголовки (Origin, Referer, User-Agent)
2. **Status Check Middleware** — проверяет HTTP статус ответа и возвращает ошибку при статусе ≥400
3. **Auth Middleware** — добавляет заголовок Authorization с access token

## Обработка ошибок

SDK возвращает структурированные ошибки с подробным описанием:

```go
feed := client.Posts.NewFeed(ctx, 10, "popular")
posts, err := feed.Next()
if err != nil {
    // Ошибка содержит информацию о коде, сообщении и HTTP статусе
    log.Printf("Ошибка при получении постов: %v", err)
    return
}
```

## Примеры

Больше примеров использования можно найти в директории [examples/](examples/):

- [examples/get_posts/](examples/get_posts/) — базовое получение постов через итератор
- [examples/hasMore/](examples/hasMore/) — пример использования HasMore() для постраничной загрузки

## Требования

- Go 1.26+

## Лицензия

MIT

## Дисклеймер

Это неофициальный SDK. Используйте на свой риск. Авторы не несут ответственности за возможные блокировки аккаунта или другие последствия использования.
