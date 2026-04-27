# Changelog

Все значимые изменения в проекте будут документированы в этом файле.

## [0.4.0] - 2026-04-27

### Добавлено

- **Модуль уведомлений**: новый пакет `api/notifications/` для работы с уведомлениями ITD API
  - `Service` - сервис для работы с уведомлениями
  - `NotificationIterator` - итератор для постраничной загрузки уведомлений
  - `NewIterator(limit)` - создание итератора уведомлений
- **Новый тип**: `types.Notification` - структура уведомления с полной информацией об акторе и целевом объекте
- **Godoc комментарии**: добавлена полная документация для всех публичных типов и методов модуля notifications

### Изменено

- **BREAKING**: удалён параметр `ctx` из конструкторов итераторов. Контекст теперь передаётся только в метод `Next(ctx)`
  - `Posts.NewFeed(ctx, tab, limit)` → `Posts.NewFeed(tab, limit)`
  - `Posts.NewUserPosts(ctx, username, limit)` → `Posts.NewUserPosts(username, limit)`
  - `Comments.NewCommentList(ctx, postID, limit)` → `Comments.NewCommentList(postID, limit)`
  - `Notifications.NewIterator(ctx, limit)` → `Notifications.NewIterator(limit)`
- **Обновлена документация**: CLAUDE.md теперь показывает правильный паттерн создания итераторов без контекста в конструкторе
- **Обновлены примеры**: все примеры использования итераторов обновлены под новый API

### Исправлено

- **Критическая ошибка**: исправлено использование `req.Body` вместо `resp.Body` в методе `getNotifications` (приводило к чтению из пустого тела запроса)

### Миграция с 0.3.1

#### Создание итераторов

**Было (0.3.1):**
```go
ctx := context.Background()
feedIter := client.Posts.NewFeed(ctx, types.FeedTabPopular, 20)
userIter := client.Posts.NewUserPosts(ctx, "username", 20)
commentIter := client.Comments.NewCommentList(ctx, postID, 20)

for feedIter.HasMore() {
    posts, err := feedIter.Next(ctx)
    // ...
}
```

**Стало (0.4.0):**
```go
// Контекст НЕ передаётся в конструктор
feedIter := client.Posts.NewFeed(types.FeedTabPopular, 20)
userIter := client.Posts.NewUserPosts("username", 20)
commentIter := client.Comments.NewCommentList(postID, 20)

// Контекст передаётся только в Next()
for feedIter.HasMore() {
    posts, err := feedIter.Next(context.Background())
    // ...
}
```

**Примечание:** Это изменение делает API более идиоматичным для Go - контекст передаётся только туда, где он реально используется (при выполнении HTTP запроса в `Next()`). Каждый вызов `Next()` может использовать свой контекст с разными таймаутами или условиями отмены.

#### Использование модуля уведомлений

**Новая функциональность (0.4.0):**
```go
// Создание итератора уведомлений
iter := client.Notifications.NewIterator(20)

for iter.HasMore() {
    notifications, err := iter.Next(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    for _, notif := range notifications {
        fmt.Printf("[%s] %s: %s\n", 
            notif.Type, 
            notif.Actor.DisplayName, 
            notif.Preview)
    }
}
```

---

## [0.3.1] - 2026-04-26

### Добавлено

- **Sentinel errors для валидации входных данных**: теперь можно проверять ошибки через `errors.Is()`
  - `ErrEmptyContent` - контент и файлы отсутствуют
  - `ErrEmptyPostID` - пустой ID поста
  - `ErrEmptyCommentID` - пустой ID комментария
  - `ErrEmptyUserID` - пустой ID пользователя
  - `ErrEmptyRefreshToken` - пустой refresh token
  - `ErrNilPoll` - poll = nil при создании поста с опросом
  - `ErrInsufficientPollOptions` - в опросе < 2 вариантов
  - `ErrEmptyReplyToUserID` - пустой ID пользователя при создании ответа
  - `ErrEmptyParentCommentID` - пустой ID родительского комментария
- **Sentinel errors для API ошибок**: маппинг HTTP статусов на типизированные ошибки
  - `ErrUnauthorized` (401) - невалидный/истёкший токен
  - `ErrNotFound` (404) - ресурс не найден
  - `ErrForbidden` (403) - недостаточно прав
  - `ErrRateLimited` (429) - превышен лимит запросов
  - `ErrServerError` (500+) - внутренняя ошибка сервера
- **Валидация файлов для комментариев**: методы `CreateComment` и `CreateReply` теперь проверяют формат и количество файлов
  - Поддерживаемые форматы: `.png`, `.webp`
  - Максимальное количество файлов: 10
- **Новые типы ошибок** в пакете `errors/`:
  - `InvalidFileExtension` - неподдерживаемое расширение файла
  - `TooManyFiles` - превышено максимальное количество файлов
  - `NoFileExtension` - файл без расширения
- **Специализированные типы для результатов создания**:
  - `CreatedPostBase` - базовая структура с общими полями
  - `CreatedPost` - результат создания простого поста
  - `CreatedPostWithPoll` - результат создания поста с опросом
  - `CreatedPostWithRepost` - результат создания репоста
  - `CreatedComment` - результат создания комментария
- **Новый пример**: `examples/comments/get_comments.go` - получение списка комментариев с использованием итератора
- **Новый пример**: `examples/errors/error_handling.go` - демонстрация всех возможностей обработки ошибок через `errors.Is()` и `errors.As()`

### Изменено

- **BREAKING**: возвращаемые типы методов создания изменены на более специфичные:
  - `Posts.Create()`: `*types.Post` → `*types.CreatedPost`
  - `Posts.CreateWithPoll()`: `*types.Post` → `*types.CreatedPostWithPoll`
  - `Posts.Repost()`: `*types.Post` → `*types.CreatedPostWithRepost`
  - `Comments.CreateComment()`: `*types.CreateComment` → `*types.CreatedComment`
  - `Comments.CreateReply()`: `*types.CreateComment` → `*types.CreatedComment`
- **Улучшена обработка ошибок**: `APIError` теперь поддерживает `Unwrap()` для проверки через `errors.Is()`
- **Реорганизация пакета `errors/`**: разделение на четыре файла для лучшей организации
  - `errors/API.go` - типы ошибок API и sentinel errors для HTTP статусов
  - `errors/transport.go` - обработка HTTP ошибок и `CheckResponse()`
  - `errors/files.go` - ошибки работы с файлами
  - `errors/validation.go` - sentinel errors для валидации входных данных
- **Рефакторинг типов Created\***: использование композиции через `CreatedPostBase` для устранения дублирования кода
- **Улучшена документация**: добавлены godoc комментарии для всех публичных типов, методов и переменных
- **Примеры**: добавлен `//go:build ignore` во все файлы примеров для корректной работы с `go run`
- **Упрощена инициализация в примерах**: использование `_ "github.com/joho/godotenv/autoload"` вместо явного вызова `godotenv.Load()`

### Удалено

- Удалён дублирующий пример `examples/posts/createPostWithFiles.go` (функциональность объединена с `createPost.go`)
- Удалён тип `types.CreateComment` (заменён на `types.CreatedComment`)

### Миграция с 0.3.0

#### Использование новых типов результатов

**Было (0.3.0):**
```go
var post *types.Post
post, err := client.Posts.Create(ctx, "Контент", "/path/to/file.jpg")
```

**Стало (0.3.1):**
```go
var post *types.CreatedPost
post, err := client.Posts.Create(ctx, "Контент", "/path/to/file.jpg")
```

**Примечание:** Новые типы содержат только необходимые поля, возвращаемые API при создании. Если вам нужна полная информация о посте, используйте метод `Get()` после создания.

#### Проверка ошибок через errors.Is()

**Новая функциональность (0.3.1):**
```go
// Валидация входных данных
post, err := client.Posts.Get(ctx, "")
if errors.Is(err, itderrors.ErrEmptyPostID) {
    log.Println("ID поста не может быть пустым")
}

// API ошибки
post, err := client.Posts.Get(ctx, "invalid-id")
if errors.Is(err, itderrors.ErrNotFound) {
    log.Println("Пост не найден")
}
if errors.Is(err, itderrors.ErrUnauthorized) {
    log.Println("Токен истёк, нужно обновить")
}

// Получение деталей API ошибки
var apiErr *itderrors.APIError
if errors.As(err, &apiErr) {
    log.Printf("API error: code=%s, status=%d", apiErr.Code, apiErr.StatusCode)
}

// Валидация файлов
post, err := client.Posts.Create(ctx, "Контент", "/path/to/file.mp3")
if errors.Is(err, itderrors.InvalidFileExtension) {
    log.Println("Неподдерживаемый формат файла")
}
```

---

## [0.3.0] - 2026-04-21

### Изменено

- **BREAKING**: Удалены интерфейсы `PostsAPI`, `UserAPI`, `CommentsAPI` из `types/interfaces.go`. Клиент теперь использует конкретные типы `posts.Service`, `user.Service`, `comments.Service` вместо интерфейсов
- **Архитектурный рефакторинг**: API слой перемещён из `internal/api/` в публичный пакет `api/`
  - `internal/api/posts/` → `api/posts/`
  - `internal/api/comments/` → `api/comments/`
  - `internal/api/user/` → `api/user/`

### Внутренние изменения

- Создан отдельный пакет `errors/` для обработки ошибок API и транспортного слоя
  - `internal/dto/errors.go` → `errors/APIerrors.go`
  - `internal/pkg/errors/httperrors.go` → `errors/transportErrors.go`
- Реорганизация итераторов:
  - Перемещён базовый итератор: `internal/pkg/iterator/` → `internal/iterator/`
  - Созданы специализированные файлы итераторов в каждом API пакете (`api/posts/iterator.go`, `api/comments/iterator.go`)
  - Удалены старые реализации итераторов (`feediterator.go`, специфичные для каждого модуля)
- Консолидация пакета аутентификации:
  - `internal/dto/auth.go` → `internal/auth/dto.go`
  - `internal/pkg/jwt/jwtparser.go` → `internal/auth/jwt.go`
- Удалены неиспользуемые файлы:
  - `internal/api/posts/getpostid.go`
  - `internal/pkg/raws/rawanswer.go` (перемещён в `internal/testutil/debug.go`)
- Обновлена документация в `CLAUDE.md` и `README.md` для отражения новой структуры проекта

### Миграция с 0.2.1

#### Использование интерфейсов

**Было (0.2.1):**
```go
var postsAPI types.PostsAPI = client.Posts
```

**Стало (0.2.2):**
```go
// Используйте конкретные типы напрямую
var postsService *posts.Service = &client.Posts
// Или просто используйте client.Posts без присваивания
```

**Примечание:** Если вы не использовали интерфейсные типы из `types/interfaces.go` напрямую в своём коде, миграция не требуется. Все публичные методы API остались без изменений.

---

## [0.2.1] - 2026-04-19

### Добавлено

- **Валидация входных данных**: методы `Create`, `CreateWithPoll`, `CreateComment`, `CreateReply` теперь проверяют корректность параметров перед отправкой запроса
- **Новый тип `AuthorInfo`**: унифицированный тип для информации об авторе контента (заменяет `Author` и `PostAuthor`)

### Изменено

- **BREAKING**: типы `Author` и `PostAuthor` объединены в единый тип `AuthorInfo`. Все поля `Author` в структурах теперь имеют тип `AuthorInfo`
- **BREAKING**: поле `IsVerified` в `PostAuthor` переименовано в `Verified` для консистентности

### Внутренние изменения

- Заменены `map[string]interface{}` на типизированные структуры в методах `Repost` и `Vote`
- Добавлен helper метод `UploadFiles` в транспортном слое для устранения дублирования кода
- Удалён неиспользуемый метод `DoJSON` из транспортного слоя
- Улучшена типобезопасность внутренних API запросов

### Миграция с 0.2.0

#### Использование типа AuthorInfo

**Было (0.2.0):**
```go
var author types.Author  // или types.PostAuthor
```

**Стало (0.2.1):**
```go
var author types.AuthorInfo
```

**Примечание:** Если вы не использовали типы `Author` или `PostAuthor` напрямую в своём коде, миграция не требуется.

---

## [0.2.0] - 2026-04-18

### Добавлено

- **Автоматическая загрузка файлов**: методы `Create`, `CreateComment`, `CreateReply` теперь принимают пути к файлам вместо ID вложений. Файлы автоматически загружаются на сервер
- **Создание постов с опросами**: новый метод `CreateWithPoll` для создания постов с прикреплёнными опросами
- **Автоматическая загрузка баннера**: метод `UpdateProfile` теперь автоматически загружает баннер при указании `BannerPath`
- **Загрузка файлов**: новый метод `Upload` в транспортном слое (`internal/transport/files.go`) для загрузки файлов на сервер ITD
- **Типы для опросов**: добавлены `PollRequest` и `PollOptionRequest` для создания опросов
- **Новый тип**: `UpdateProfile` - структура для обновления профиля с автозагрузкой баннера
- **Примеры использования**:
  - `examples/posts/createPostWithFiles.go` - создание поста с файлами
  - `examples/posts/createPostWithPoll.go` - создание поста с опросом

### Изменено

- **BREAKING**: сигнатура метода `Posts.Create` изменена с `Create(ctx, content, attachmentIDs...)` на `Create(ctx, content, filePaths...)`. Теперь метод принимает пути к файлам и автоматически загружает их
- **BREAKING**: сигнатура метода `Comments.CreateComment` изменена с `CreateComment(ctx, postID, content, attachmentIDs)` на `CreateComment(ctx, postID, content, filePaths...)`. Теперь метод принимает пути к файлам и автоматически загружает их
- **BREAKING**: сигнатура метода `Comments.CreateReply` изменена с `CreateReply(ctx, commentID, content, attachmentIDs)` на `CreateReply(ctx, parentCommentID, replyToUserID, content, filePaths...)`. Теперь метод принимает пути к файлам и автоматически загружает их
- **BREAKING**: сигнатура метода `User.UpdateProfile` изменена с `UpdateProfile(ctx, *displayName, *username, *bio, *bannerID)` на `UpdateProfile(ctx, config)`. Теперь метод принимает структуру `UpdateProfile` и автоматически загружает баннер
- **BREAKING**: метод `CreateWithPoll` принимает указатель на `PollRequest` (`*PollRequest`) вместо значения
- **Переименование типов**: `PollOption` переименован в `PollOptionResponse` для ясности (отличие от `PollOptionRequest`)
- Обновлены примеры `examples/posts/createPost.go`, `examples/comments/create_comment.go`, `examples/comments/create_reply.go`, `examples/user/updateProfile.go` с комментариями о новой функциональности

### Внутренние изменения

- Приватизированы внутренние DTO типы: `ResponseFeed` → `responseFeed`, `ResponsePost` → `responsePost`
- Добавлена структура `createPostRequest` для формирования запросов на создание постов
- Добавлена структура `updateCommentRequest` для формирования запросов на обновление комментариев
- Добавлен метод `NewRequestMultipart` в транспортном слое для multipart/form-data запросов
- Версия SDK обновлена с `0.1.0` до `0.2.0`

### Миграция с 0.1.x

#### Создание постов с файлами

**Было (0.1.x):**
```go
// Нужно было сначала загрузить файлы и получить их ID
attachmentID := "..." // ID загруженного файла
post, err := client.Posts.Create(ctx, "Контент", attachmentID)
```

**Стало (0.2.0):**
```go
// Просто передайте пути к файлам - они загрузятся автоматически
post, err := client.Posts.Create(ctx, "Контент", "/path/to/file.jpg")
```

#### Создание комментариев с файлами

**Было (0.1.x):**
```go
comment, err := client.Comments.CreateComment(ctx, postID, "Текст", []string{attachmentID})
```

**Стало (0.2.0):**
```go
comment, err := client.Comments.CreateComment(ctx, postID, "Текст", "/path/to/file.jpg")
```

#### Обновление профиля

**Было (0.1.x):**
```go
displayName := "Новое имя"
bio := "Новая биография"
bannerID := "banner_id"
user, err := client.User.UpdateProfile(ctx, &displayName, nil, &bio, &bannerID)
```

**Стало (0.2.0):**
```go
profile := types.UpdateProfile{
    DisplayName: "Новое имя",
    Bio:         "Новая биография",
    BannerPath:  "/path/to/banner.jpg", // Автоматически загрузится
}
user, err := client.User.UpdateProfile(ctx, profile)
```

#### Создание постов с опросами

**Новая функциональность (0.2.0):**
```go
poll := types.PollRequest{
    Question: "Ваш любимый язык?",
    Options: []types.PollOptionRequest{
        {Text: "Go"},
        {Text: "Python"},
    },
    MultipleChoice: false,
}
post, err := client.Posts.CreateWithPoll(ctx, "Опрос!", &poll)
```

---

## [0.1.0] - 2026-03-XX

### Добавлено

- Первый публичный релиз ITD Go SDK
- Поддержка основных операций с постами (создание, удаление, лайки, репосты)
- Поддержка работы с комментариями
- Поддержка работы с пользователями
- Итераторы для пагинации постов и комментариев
- Автоматическое управление токенами аутентификации
- Примеры использования для всех основных методов
