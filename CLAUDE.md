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
│  - Posts, User, Comments            │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  API Layer (api/)                   │
│  - posts, user, comments             │
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
// Универсальный интерфейс
type Iterator[T any] interface {
    HasMore() bool
    Next() ([]T, error)
}

// Конкретные типы
type FeedIterator = Iterator[*Post]
type CommentIterator = Iterator[*Comment]
```

**Важно:** Все итераторы возвращают **указатели** на элементы (`*Post`, `*Comment`), а не значения.

### 2. Аутентификация

SDK использует refresh token из cookies браузера:

1. Пользователь передаёт refresh token в `Config`
2. SDK автоматически получает access token через `/api/v1/auth/refresh`
3. Access token добавляется к каждому запросу через middleware
4. При истечении токена происходит автоматическое обновление

### 3. Автоматическая загрузка файлов

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

### 4. Типы и интерфейсы

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

// YourIterator предоставляет интерфейс для постраничной загрузки данных.
type YourIterator interface {
    // HasMore возвращает true, если есть ещё данные для загрузки.
    HasMore() bool
    // Next загружает и возвращает следующую страницу данных.
    // Параметры:
    //   - ctx: контекст для управления временем жизни запроса
    Next(ctx context.Context) ([]*types.YourType, error)
}
```

### 2. Создать функцию-конструктор

В том же файле `internal/api/yourmodule/iterator.go`:

```go
func newYourIterator(s *Service, limit int) YourIterator {
    fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.YourType, iterator.PageToken, bool, error) {
        result, err := s.getYourData(ctx, token.Cursor, limit)
        if err != nil {
            return nil, iterator.PageToken{}, false, err
        }
        next := iterator.PageToken{Cursor: result.NextCursor}
        return result.Items, next, result.HasMore, nil
    }
    
    return iterator.New[*types.YourType](fetch, iterator.PageToken{})
}
```

### 3. Добавить публичный метод

В файле `api/yourmodule/yourmodule.go`:

```go
func (s *Service) NewYourIterator(limit int) YourIterator {
    return newYourIterator(s, limit)
}
```

### 4. Использование итератора

```go
iter := service.NewYourIterator(20)

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
