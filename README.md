<div align="center">

<img width="400" alt="Image" src="https://github.com/user-attachments/assets/8517c4b6-fd5f-498d-befc-8896976453dd" />


**Неофициальный Go SDK для работы с API социальной сети** [итд.com](https://итд.com)

**Версия:** 0.2.1

Made with ❤️ by [@Slipynil](https://github.com/slipynil)

[![Go Reference](https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/slipynil/itd-go)
[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/slipynil_chan)

</div>

## Текущие реализованные методы

### Posts API
```go
// Итераторы для постраничной загрузки
feed := client.Posts.NewFeed(ctx, types.FeedTabPopular, 20)           // Лента постов
userPosts := client.Posts.NewUserPosts(ctx, "username", 20)            // Посты пользователя

// Работа с постами
post, err := client.Posts.Get(ctx, postID)                             // Получение поста
post, err := client.Posts.Create(ctx, "Текст")                         // Создание поста
post, err := client.Posts.Create(ctx, "Текст", "/path/to/image.jpg")  // С автозагрузкой файлов
post, err := client.Posts.CreateWithPoll(ctx, "Текст", &poll)           // Создание поста с опросом
post, err := client.Posts.CreateWithPoll(ctx, "Текст", &poll, "/path/to/file.jpg") // С опросом и файлами
post, err := client.Posts.Repost(ctx, postID, "Комментарий")           // Репост
err := client.Posts.Delete(ctx, postID)                                // Удаление поста

// Взаимодействие
result, err := client.Posts.Like(ctx, postID)                          // Лайк
result, err := client.Posts.Unlike(ctx, postID)                        // Убрать лайк
poll, err := client.Posts.Vote(ctx, postID, optionIDs...)              // Голосование в опросе
view, err := client.Posts.View(ctx, postID)                            // Отметить просмотр
```

### User API
```go
user, err := client.User.Me(ctx)                                       // Текущий пользователь
user, err := client.User.Get(ctx, "username")                          // Получение пользователя
err := client.User.Follow(ctx, "username")                             // Подписаться
err := client.User.Unfollow(ctx, "username")                           // Отписаться
followers, err := client.User.GetFollowers(ctx, "username", 20)        // Список подписчиков
following, err := client.User.GetFollowing(ctx, "username", 20)        // Список подписок

// Обновление профиля с автозагрузкой баннера
profile := types.UpdateProfile{
    DisplayName: "Новое имя",
    Bio:         "Новая биография",
    BannerPath:  "/path/to/banner.jpg", // Автоматически загрузится
}
response, err := client.User.UpdateProfile(ctx, profile)
```

### Comments API
```go
// Итератор для комментариев
comments := client.Comments.NewCommentList(ctx, postID, 20)            // Комментарии к посту
replies, err := client.Comments.ListReplies(ctx, commentID, 20)        // Ответы на комментарий

// Создание
comment, err := client.Comments.CreateComment(ctx, postID, "Текст")    // Создать комментарий
comment, err := client.Comments.CreateComment(ctx, postID, "Текст", "/path/to/file.jpg") // С файлом
reply, err := client.Comments.CreateReply(ctx, parentCommentID, replyToUserID, "Текст") // Ответить
reply, err := client.Comments.CreateReply(ctx, parentCommentID, replyToUserID, "Текст", "/path/to/file.jpg") // С файлом

// Взаимодействие
err := client.Comments.Like(ctx, commentID)                            // Лайк
err := client.Comments.Unlike(ctx, commentID)                          // Убрать лайк
updated, err := client.Comments.Update(ctx, commentID, "Новый текст") // Обновить
err := client.Comments.Delete(ctx, commentID)                          // Удалить
```
## Установка

```bash
go get github.com/slipynil/itd-go
```

## Быстрый старт

```go
package main

import (
    "context"
    "fmt"
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
    feed := client.Posts.NewFeed(ctx, types.FeedTabPopular, 20)

    // Получить первую страницу постов
    for feed.HasMore() {
        posts, err := feed.Next()
        if err != nil {
            log.Fatal(err)
        }

        for _, post := range posts {
            fmt.Printf("%s: %s\n", post.Author.DisplayName, post.Content)
        }
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

- `NewFeed(ctx, tab, limit)` — итератор для ленты постов (tab: types.FeedTabPopular, types.FeedTabClan, types.FeedTabFollowing)
- `NewUserPosts(ctx, username, limit)` — итератор для постов пользователя
- `Get(ctx, postID)` — получение поста по ID
- `Create(ctx, content, filePaths...)` — создание поста с автоматической загрузкой файлов (filePaths опциональны)
- `CreateWithPoll(ctx, content, *poll, filePaths...)` — создание поста с опросом и файлами (poll - указатель на PollRequest, filePaths опциональны)
- `Delete(ctx, postID)` — удаление поста
- `Like(ctx, postID)` — лайк на пост
- `Unlike(ctx, postID)` — убрать лайк
- `Repost(ctx, postID, content)` — репост с опциональным комментарием
- `Vote(ctx, postID, optionIDs...)` — голосование в опросе
- `View(ctx, postID)` — отметить пост как просмотренный

### User

- `Me(ctx)` — получение информации о текущем пользователе
- `Get(ctx, username)` — получение информации о пользователе по username
- `Follow(ctx, username)` — подписаться на пользователя
- `Unfollow(ctx, username)` — отписаться от пользователя
- `GetFollowers(ctx, username, limit)` — получить список подписчиков пользователя
- `GetFollowing(ctx, username, limit)` — получить список подписок пользователя
- `UpdateProfile(ctx, config)` — обновить профиль (config - структура UpdateProfile с полями для обновления, пустые поля не изменяются, BannerPath автоматически загружается)

### Comments

- `NewCommentList(ctx, postID, limit)` — итератор для комментариев к посту
- `ListReplies(ctx, commentID, limit)` — получить список ответов на комментарий
- `CreateComment(ctx, postID, content, filePaths...)` — создать комментарий с автоматической загрузкой файлов (filePaths опциональны)
- `CreateReply(ctx, parentCommentID, replyToUserID, content, filePaths...)` — создать ответ с автоматической загрузкой файлов (filePaths опциональны)
- `Delete(ctx, commentID)` — удалить комментарий
- `Like(ctx, commentID)` — лайк на комментарий
- `Unlike(ctx, commentID)` — убрать лайк с комментария
- `Update(ctx, commentID, content)` — обновить содержимое комментария

**Примечание:** Итератор `NewCommentList` реализует интерфейс `CommentIterator` с методами `HasMore()` и `Next()`.

### Итераторы

Все итераторы (`NewFeed`, `NewUserPosts`) реализуют интерфейс `FeedIterator`:
- `HasMore()` — проверка наличия следующей страницы
- `Next()` — загрузка следующей страницы постов

## Архитектура

SDK построен на многоуровневой архитектуре с чёткими границами ответственности:

```
itd-go/
├── client.go           — публичная точка входа
├── config.go           — публичная конфигурация
│
├── types/              — публичные интерфейсы и типы данных
│   ├── interfaces.go   — API интерфейсы (PostsAPI, UserAPI, CommentsAPI)
│   ├── post.go         — структуры постов
│   ├── user.go         — структуры пользователей
│   └── comment.go      — структуры комментариев
│
└── internal/
    ├── root/           — создание HTTP клиента, установка заголовков
    ├── auth/           — управление токенами, аутентификация
    ├── transport/      — HTTP механика, middleware
    └── api/            — группы методов (posts, user, comments)
        └── posts/
            ├── posts.go        — публичные методы API
            ├── iterator.go     — generic iterator для пагинации
            ├── feediterator.go — фабрики итераторов
            ├── getpostid.go    — внутренние методы запросов
            └── dto.go          — внутренние структуры ответов
```

### Generic Iterator Pattern

Для устранения дублирования кода используется generic iterator:

```go
type Iterator[T any] interface {
    HasMore() bool
    Next() ([]T, error)
}
```

Это позволяет переиспользовать логику пагинации для разных типов данных (посты, комментарии и т.д.).

### Поток запроса

```
client.Posts.NewFeed(ctx, limit, sort)
    → newFeedIterator()
        → newPaginatedIterator[types.Post](ctx, fetchFunc)
    → iterator.Next()
        → fetchFunc(ctx, cursor)
            → getFeed(ctx, limit, cursor, sort)
                → transport.NewRequest()
                → transport.Do()
                    → Middleware цепочка:
                        1. Root Transport (Origin, Referer, User-Agent)
                        2. Status Check (проверка статуса ≥400)
                        3. Auth Middleware (Authorization токен)
                    → HTTP запрос
                → JSON декодирование
            → возврат (posts, nextCursor, hasMore)
        → обновление состояния итератора
        → возврат постов

Аналогично для NewUserPosts() с getUserPosts() вместо getFeed()
```

### Middleware цепочка

1. **Root Transport** — устанавливает базовые заголовки (Origin, Referer, User-Agent)
2. **Status Check Middleware** — проверяет HTTP статус ответа и возвращает ошибку при статусе ≥400
3. **Auth Middleware** — добавляет заголовок Authorization с access token

## Обработка ошибок

SDK возвращает структурированные ошибки с подробным описанием:

```go
feed := client.Posts.NewFeed(ctx, types.FeedTabPopular, 10)
posts, err := feed.Next()
if err != nil {
    // Ошибка содержит информацию о коде, сообщении и HTTP статусе
    log.Printf("Ошибка при получении постов: %v", err)
    return
}
```

## Примеры

Примеры использования находятся в директории `examples/`:

### Posts (`examples/posts/`)

- `showFeed.go` — получение ленты постов через итератор
- `showUserPosts.go` — получение постов пользователя через итератор
- `getPost.go` — получение поста по ID
- `createPost.go` — создание нового поста
- `createPostWithFiles.go` — создание поста с автоматической загрузкой файлов
- `createPostWithPoll.go` — создание поста с опросом
- `deletePost.go` — удаление поста
- `like.go` — лайк на пост
- `unlike.go` — убрать лайк
- `repost.go` — репост поста
- `vote.go` — голосование в опросе

### User (`examples/user/`)

- `meInfo.go` — получение информации о текущем пользователе
- `userInfo.go` — получение информации о пользователе по username
- `follow.go` — подписка на пользователя
- `unfollow.go` — отписка от пользователя
- `get_followers.go` — получение списка подписчиков
- `get_following.go` — получение списка подписок
- `updateProfile.go` — обновление профиля пользователя

### Comments (`examples/comments/`)

- `get_comment_list.go` — получение списка комментариев через итератор
- `get_one_comment.go` — получение одного комментария
- `get_reply_list.go` — получение списка ответов на комментарий
- `create_comment.go` — создание комментария к посту
- `create_reply.go` — создание ответа на комментарий
- `update_comment.go` — обновление содержимого комментария
- `delete_comment.go` — удаление комментария
- `like_comment.go` — лайк на комментарий
- `unlike_comment.go` — убрать лайк с комментария

## Требования

- Go 1.26+

## Лицензия

MIT

## Дисклеймер

Это неофициальный SDK. Используйте на свой риск. Авторы не несут ответственности за возможные блокировки аккаунта или другие последствия использования.
