# Changelog

Все значимые изменения в проекте будут документированы в этом файле.

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
