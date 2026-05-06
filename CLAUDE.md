# ITD Go SDK - Документация для разработки

## О проекте

**itd-go** — неофициальный Go SDK для работы с API социальной сети [итд.com](https://итд.com).

Версия: 0.2.0  
Язык: Go 1.26+

## Архитектура

Проект использует многоуровневую архитектуру:

```
┌─────────────────────────────────────┐
│  Client (публичный API)             │
│  - Posts, User, Comments,           │
│    Notifications                    │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  API Layer (api/)                   │
│  - posts, user, comments,           │
│    notifications                    │
│  - Бизнес-логика, итераторы        │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Transport (internal/transport/)    │
│  - HTTP клиент, middleware          │
│  - Сериализация/десериализация      │
│  - Загрузка файлов (files.go)       │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Auth (internal/auth/)              │
│  - Управление токенами              │
│  - Refresh token → Access token     │
└─────────────────────────────────────┘
```

## Ключевые концепции

### 1. Итераторы для пагинации

Все методы, возвращающие списки, используют паттерн Iterator:

```go
// Каждый API модуль определяет свой интерфейс Iterator
// api/posts/iterator.go
type Iterator interface {
    HasMore() bool
    Next(ctx context.Context) ([]*types.Post, error)
}

// api/comments/iterator.go
type Iterator interface {
    HasMore() bool
    Next(ctx context.Context) ([]*types.Comment, error)
}

// api/notifications/iterator.go
type Iterator interface {
    HasMore() bool
    Next(ctx context.Context) ([]*types.Notification, error)
}
```

**Важно:** 
- Все итераторы возвращают **указатели** на элементы (`*Post`, `*Comment`, `*Notification`), а не значения
- Интерфейс `Iterator` определяется локально в каждом пакете (Go best practice)
- Контекст передаётся только в метод `Next(ctx)`, не в конструктор


### 2. Аутентификация

SDK использует refresh token из cookies браузера:

1. Пользователь передаёт refresh token в `Config`
2. SDK автоматически получает access token через `/api/v1/auth/refresh`
3. Access token добавляется к каждому запросу через middleware
4. При истечении токена происходит автоматическое обновление

### 3. Автоматическая обработка rate limiting (429)

SDK автоматически повторяет запросы при получении ошибки 429 (Too Many Requests):

1. При ошибке 429 SDK автоматически повторяет запрос
2. Используется exponential backoff: 1s → 2s → 4s → 8s
3. По умолчанию 3 попытки с начальной задержкой 1 секунда
4. Можно настроить через `Config.MaxRetries` и `Config.RetryDelay`

```go
cfg := itdgo.Config{
    RefreshToken: "your_token",
    MaxRetries:   5,                    // 5 попыток вместо 3
    RetryDelay:   2 * time.Second,      // начальная задержка 2 секунды
}

// Для отключения retry логики установите MaxRetries = 0
cfg := itdgo.Config{
    RefreshToken: "your_token",
    MaxRetries:   0,  // отключить автоматические повторы
}
```

**Важно:** Retry применяется только к ошибкам 429. Другие ошибки (401, 404, 5xx) возвращаются немедленно без повторов.

### 4. Автоматическая загрузка файлов

SDK автоматически загружает файлы на сервер при создании постов и комментариев:

1. Методы принимают пути к файлам (`filePaths ...string`)
2. SDK загружает каждый файл через `transport.Upload(ctx, path)`
3. Получает ID загруженных файлов
4. Автоматически прикрепляет их к посту/комментарию

```go
// Пользователь просто передаёт пути к файлам
post, err := client.Posts.Create(ctx, "Контент", "/path/to/image.jpg")

// SDK автоматически:
// 1. Загружает файл через transport.Upload()
// 2. Получает attachment ID
// 3. Создаёт пост с этим ID
```

**Важно:** Методы `Create`, `CreateWithPoll`, `CreateComment`, `CreateReply` принимают `filePaths ...string`, а не `attachmentIDs`.

### 5. Десериализация дат с transport.DataOptions

**Критически важно:** Все методы, которые десериализуют типы с полями `time.Time`, **обязаны** использовать `transport.DataOptions`.

`DataOptions` содержит кастомный unmarshaler для `time.Time`, который поддерживает несколько форматов дат от API:
- RFC3339 (ISO8601): `2024-01-15T10:30:00Z`
- С микросекундами: `2024-01-15 10:30:00.123456+03`
- Без микросекунд: `2024-01-15 10:30:00+03`

```go
// ✅ ПРАВИЛЬНО - используется DataOptions
var result types.Post
if err := json.UnmarshalRead(resp.Body, &result, transport.DataOptions); err != nil {
    return nil, err
}

// ❌ НЕПРАВИЛЬНО - даты могут распарситься некорректно
var result types.Post
if err := json.UnmarshalRead(resp.Body, &result); err != nil {
    return nil, err
}
```

**Типы, требующие DataOptions:**
- `types.Post` (CreatedAt, EditedAt)
- `types.CreatedPost*` (CreatedAt)
- `types.Comment` (CreatedAt)
- `types.CreatedComment` (CreatedAt)
- `types.CommentUpdate` (UpdatedAt)
- `types.User`, `types.Me` (CreatedAt)
- `types.UpdateProfileResponse` (UpdatedAt)
- `types.Notification` (CreatedAt, ReadAt)
- `types.StreamNotification` (CreatedAt, ReadAt)

**Типы, НЕ требующие DataOptions:**
- `types.Attachment` (нет полей с датами)
- `types.UserCompact` (нет полей с датами)
- `types.Hashtag`, `types.SearchResult` (нет полей с датами)

### 6. Типы и интерфейсы

- **types/** — публичные типы и интерфейсы API
- **internal/dto/** — внутренние DTO для парсинга ответов API
- **api/*/dto.go** — DTO специфичные для каждого API модуля

## Стандарты кодирования

### Комментарии (godoc)

**Обязательно** для всех публичных типов, полей и методов:

```go
// Post представляет пост в социальной сети ITD.
// Содержит текстовый контент, метаданные и статистику.
type Post struct {
    // ID - уникальный идентификатор поста
    ID string `json:"id"`
    
    // Content - текстовое содержимое поста
    Content string `json:"content"`
}
```

### Указатели vs значения

- **Итераторы:** всегда возвращают указатели (`[]*Post`, `[]*Comment`)
- **Рекурсивные структуры:** обязательно используют указатели (`Replies []*Comment`)
- **API методы:** возвращают указатели для больших структур (`*Post`, `*User`)

### Обработка ошибок

```go
// Все ошибки API оборачиваются в APIError
type APIError struct {
    Code       string  // код ошибки от API
    Message    string  // описание
    StatusCode int     // HTTP статус
}
```

## Как добавить новый API метод

### 1. Добавить DTO (если нужно)

В `api/posts/dto.go`:

```go
// ResponseNewMethod представляет ответ API для нового метода.
type ResponseNewMethod struct {
    // Data - данные ответа
    Data DataStruct `json:"data"`
}
```

### 2. Реализовать метод

В `api/posts/posts.go`:

```go
// NewMethod описание метода.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - param: описание параметра
//
// Возвращает результат или ошибку при проблемах с сетью/API.
func (s *Service) NewMethod(ctx context.Context, param string) (*types.Result, error) {
    path := fmt.Sprintf("/api/posts/%s", param)
    
    req, err := p.transport.NewRequest(ctx, "GET", path, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := p.transport.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result ResponseNewMethod
    if err := json.UnmarshalRead(resp.Body, &result, transport.DataOptions); err != nil {
        return nil, err
    }
    
    return &result.Data, nil
}
```

### 3. Добавить пример использования

В `examples/posts/new_method.go`:

```go
package main

import (
    "context"
    "log"
    
    itdgo "github.com/slipynil/itd-go"
)

func main() {
    cfg := itdgo.Config{
        RefreshToken: "your_token",
    }
    
    client, err := itdgo.New(context.Background(), cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    result, err := client.Posts.NewMethod(context.Background(), "param")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println(result)
}
```

## Как добавить метод с загрузкой файлов

Если метод должен поддерживать загрузку файлов:

### 1. Реализовать с автозагрузкой

В `api/posts/posts.go`:

```go
// CreateSomething создаёт что-то с файлами.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - content: текстовое содержимое
//   - filePaths: пути к файлам для автоматической загрузки
//
// Возвращает результат или ошибку при проблемах с сетью/API.
func (s *Service) CreateSomething(ctx context.Context, content string, filePaths ...string) (*types.Result, error) {
    // Загружаем файлы и получаем их ID
    attachmentIDs := make([]string, 0, len(filePaths))
    for _, path := range filePaths {
        attachment, err := p.transport.Upload(ctx, path)
        if err != nil {
            return nil, err
        }
        attachmentIDs = append(attachmentIDs, attachment.ID)
    }
    
    // Создаём payload с ID загруженных файлов
    payload := createRequest{
        Content:       content,
        AttachmentIDs: attachmentIDs,
    }
    
    // Отправляем запрос
    req, err := p.transport.NewRequest(ctx, "POST", "/api/something", payload)
    if err != nil {
        return nil, err
    }
    
    resp, err := p.transport.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result types.Result
    if err := json.UnmarshalRead(resp.Body, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

## API уведомлений

Модуль `api/notifications` предоставляет методы для работы с уведомлениями.

### Основные методы

#### ListUnread - получение непрочитанных уведомлений

```go
// Получить все непрочитанные уведомления
notifications, err := client.Notifications.ListUnread(ctx)
if err != nil {
    log.Fatal(err)
}

for _, n := range notifications {
    fmt.Printf("[%s] %s: %s\n", n.Type, n.Actor.DisplayName, n.Preview)
}
```

**Важно:** Метод останавливается при первом прочитанном уведомлении, так как API всегда возвращает уведомления в хронологическом порядке (новые первыми).

#### MarkRead - пометка по ID

```go
// Пометить конкретные уведомления как прочитанные
err := client.Notifications.MarkRead(ctx, "id1", "id2", "id3")
if err != nil {
    log.Fatal(err)
}
```

#### MarkNotificationsRead - пометка загруженных уведомлений

```go
// Получить непрочитанные
notifications, err := client.Notifications.ListUnread(ctx)
if err != nil {
    log.Fatal(err)
}

// Показать пользователю
for _, n := range notifications {
    fmt.Println(n.Preview)
}

// Пометить как прочитанные (без повторной загрузки)
err = client.Notifications.MarkNotificationsRead(ctx, notifications)
if err != nil {
    log.Fatal(err)
}
```

**Преимущество:** Избегает повторной загрузки уведомлений, если они уже получены.

#### MarkAllRead - пометка всех непрочитанных

```go
// Пометить все непрочитанные уведомления одной командой
err := client.Notifications.MarkAllRead(ctx)
if err != nil {
    log.Fatal(err)
}
```

**Примечание:** Метод внутри вызывает `ListUnread()` и затем `MarkRead()`. Если уведомления уже загружены, используйте `MarkNotificationsRead()`.

#### Stream - получение уведомлений в реальном времени

```go
// Открыть SSE стрим для получения уведомлений
stream, errs := client.Notifications.Stream(ctx)

for {
    select {
    case n, ok := <-stream:
        if !ok {
            return // стрим закрыт
        }
        fmt.Printf("[%s] %s: %s\n", n.Type, n.Actor.DisplayName, n.Preview)
        
    case err := <-errs:
        log.Printf("Stream error: %v", err)
        return // при ошибке стрим автоматически закрывается
    }
}
```

**Важно:** 
- Метод возвращает два канала: `<-chan *types.StreamNotification` для данных и `<-chan error` для ошибок
- При получении ошибки стрим автоматически закрывается
- Автоматическое переподключение не реализовано - нужно реализовывать на стороне клиента
- Стрим пропускает heartbeat события (keepalive от сервера без ID)
- `StreamNotification` отличается от `Notification`: содержит дополнительные поля `UserID`, `Sound` и nullable `ReadAt`

### Выбор метода

- **`MarkRead(ctx, ids...)`** - когда есть только ID уведомлений
- **`MarkNotificationsRead(ctx, notifications)`** - когда уведомления уже загружены
- **`MarkAllRead(ctx)`** - когда нужно пометить всё одной командой без предварительной загрузки
- **`Stream(ctx)`** - когда нужно получать уведомления в реальном времени через SSE

### TODO: Улучшения для Stream()

**Кастомный логгер для обработки ошибок парсинга**

Текущее поведение: при ошибке парсинга JSON весь стрим закрывается (fail-fast подход).

Планируется добавить возможность передавать кастомный логгер в конфигурацию клиента:
```go
type Config struct {
    // ...
    Logger Logger // интерфейс для логирования ошибок стрима
}

type Logger interface {
    Error(msg string, err error)
    Warn(msg string)
}
```

Это позволит:
- Логировать ошибки парсинга без закрытия стрима (resilient режим)
- Пользователю самому решать, как обрабатывать невалидные события
- Отлаживать проблемы с API без потери соединения

**Причина:** Если API итд.com иногда отправляет невалидные события, текущий подход приводит к частым переподключениям. Кастомный логгер позволит пропускать битые события и продолжать работу стрима.

## Как создать новый итератор

### 1. Определить интерфейс в пакете, где он используется

**Важно:** В Go интерфейсы определяются там, где они используются, а не в центральном пакете types.

В файле `api/yourmodule/iterator.go`:

```go
package yourmodule

import (
    "context"
    "github.com/slipynil/itd-go/types"
)

// Iterator предоставляет интерфейс для постраничной загрузки данных.
type Iterator interface {
    // HasMore возвращает true, если есть ещё данные для загрузки.
    HasMore() bool
    // Next загружает и возвращает следующую страницу данных.
    // Параметры:
    //   - ctx: контекст для управления временем жизни запроса
    Next(ctx context.Context) ([]*types.YourType, error)
}
```

### 2. Создать функцию-конструктор

В том же файле `api/yourmodule/iterator.go`:

```go
// newYourData создаёт итератор для получения ваших данных.
// Имя функции должно описывать ЧТО итерируется.
func newYourData(s *Service, limit int) Iterator {
    fetch := func(ctx context.Context, token *iterator.PageToken) ([]*types.YourType, *iterator.PageToken, bool, error) {
        cursor := ""
        if token != nil {
            cursor = token.Cursor
        }

        result, err := s.getYourData(ctx, cursor, limit)
        if err != nil {
            return nil, nil, false, err
        }

        var next *iterator.PageToken
        if result.HasMore {
            next = &iterator.PageToken{Cursor: result.NextCursor}
        }

        return result.Items, next, result.HasMore, nil
    }
    
    return iterator.New[*types.YourType](fetch, nil)
}
```

### 3. Добавить публичный метод

В файле `api/yourmodule/yourmodule.go`:

```go
// NewYourData создаёт итератор для получения ваших данных.
// Имя метода должно описывать ЧТО итерируется.
func (s *Service) NewYourData(limit int) Iterator {
    return newYourData(s, limit)
}
```

### 4. Использование итератора

```go
iter := service.NewYourData(20)

for iter.HasMore() {
    items, err := iter.Next(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    // обработка items
}
```

**Важно:** 
- Контекст передаётся только в метод `Next(ctx)` при каждом вызове
- Не передавайте контекст в конструктор итератора - он там не используется
- Не храните контекст в структуре итератора (антипатерн в Go)
- `PageToken` используется как указатель: `nil` = первый запрос, `&PageToken{...}` = есть cursor
- Интерфейс всегда называется просто `Iterator` (в каждом пакете свой)
- Имена методов и функций должны описывать ЧТО итерируется (NewHashtagPosts, NewPostComments, NewNotifications)
```

## Структура проекта

```
itd-go/
├── client.go              # Главный клиент SDK
├── config.go              # Конфигурация SDK
├── types/                 # Публичные типы и интерфейсы
│   ├── post.go           # Типы для постов
│   ├── comment.go        # Типы для комментариев
│   ├── user.go           # Типы для пользователей
│   ├── feedTab.go        # Enum для типов ленты
│   └── pin.go            # Типы для значков
├── api/                  # Публичные API реализации
│   ├── posts/            # API постов
│   ├── comments/         # API комментариев
│   ├── notifications/    # API уведомлений
│   └── user/             # API пользователей
├── internal/
│   ├── auth/             # Аутентификация
│   ├── transport/        # HTTP транспорт
│   │   ├── transport.go  # HTTP клиент и middleware
│   │   └── files.go      # Загрузка файлов на сервер
│   ├── root/             # Корневой клиент
│   ├── dto/              # Общие DTO
│   └── pkg/              # Утилиты
│       ├── errors/       # Обработка ошибок
│       ├── iterator/     # Базовый итератор
│       ├── jwt/          # JWT парсер
│       └── raws/         # Raw JSON обработка
└── examples/             # Примеры использования
    ├── posts/
    ├── comments/
    ├── notifications/
    └── user/
```

## Важные файлы

- **client.go** — точка входа в SDK
- **internal/pkg/iterator/iterator.go** — базовая реализация итератора
- **internal/transport/transport.go** — HTTP клиент
- **internal/transport/files.go** — загрузка файлов на сервер
- **internal/auth/auth.go** — управление токенами

## Тестирование

Для тестирования нужен валидный refresh token:

1. Создать `.env` файл в `examples/`:
   ```env
   REFRESH_TOKEN=your_refresh_token_here
   USER_AGENT=Mozilla/5.0...
   ```

2. Запустить пример:
   ```bash
   cd examples && go run posts/showFeed.go
   ```

## Получение refresh token

1. Открыть [итд.com](https://итд.com) в браузере
2. Войти в аккаунт
3. Открыть DevTools → Application → Cookies
4. Скопировать значение cookie `refresh_token`

## Известные ограничения

- API не документировано официально
- Некоторые методы могут измениться без предупреждения
- Rate limiting не реализован (используйте разумные задержки)
- Нет поддержки WebSocket для real-time обновлений

## Соглашения о коммитах

- `feat:` — новая функциональность
- `fix:` — исправление бага
- `docs:` — изменения в документации
- `refactor:` — рефакторинг без изменения функциональности
- `test:` — добавление тестов
- `chore:` — обновление зависимостей, конфигурации

## Полезные команды

```bash
# Компиляция
go build ./...

# Проверка кода
go vet ./...

# Форматирование
go fmt ./...

# Запуск примера
cd examples && go run posts/showFeed.go
```

## Контакты и поддержка

Проект находится в активной разработке. При обнаружении проблем создавайте issue в репозитории.
