# ITD Go SDK - Документация для разработки

## О проекте

**itd-go** — неофициальный Go SDK для работы с API социальной сети [итд.com](https://итд.com).

Версия: 0.1  
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
│  API Layer (internal/api/)          │
│  - posts, user, comments            │
│  - Бизнес-логика, итераторы        │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Transport (internal/transport/)    │
│  - HTTP клиент, middleware          │
│  - Сериализация/десериализация      │
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

### 3. Типы и интерфейсы

- **types/** — публичные типы и интерфейсы API
- **internal/dto/** — внутренние DTO для парсинга ответов API
- **internal/api/*/dto.go** — DTO специфичные для каждого API модуля

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

### 1. Определить интерфейс в types/interfaces.go

```go
type PostsAPI interface {
    // NewMethod описание метода.
    // Параметры:
    //   - ctx: контекст
    //   - param: описание параметра
    // Возвращает результат или ошибку.
    NewMethod(ctx context.Context, param string) (*Result, error)
}
```

### 2. Добавить DTO (если нужно)

В `internal/api/posts/dto.go`:

```go
// ResponseNewMethod представляет ответ API для нового метода.
type ResponseNewMethod struct {
    // Data - данные ответа
    Data DataStruct `json:"data"`
}
```

### 3. Реализовать метод

В `internal/api/posts/posts.go`:

```go
// NewMethod реализует PostsAPI.NewMethod.
func (p *Posts) NewMethod(ctx context.Context, param string) (*types.Result, error) {
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

### 4. Добавить пример использования

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

## Как создать новый итератор

### 1. Определить тип в types/interfaces.go

```go
type NewIterator = Iterator[*NewType]
```

### 2. Создать функцию-конструктор

```go
func newNewIterator(client *Client, ctx context.Context, limit int) types.NewIterator {
    fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.NewType, iterator.PageToken, bool, error) {
        result, err := client.getNewData(ctx, token.Cursor, limit)
        if err != nil {
            return nil, iterator.PageToken{}, false, err
        }
        next := iterator.PageToken{Cursor: result.NextCursor}
        return result.Items, next, result.HasMore, nil
    }
    
    return iterator.New[*types.NewType](ctx, fetch, iterator.PageToken{})
}
```

### 3. Добавить публичный метод

```go
func (c *Client) NewIterator(ctx context.Context, limit int) types.NewIterator {
    return newNewIterator(c, ctx, limit)
}
```

## Структура проекта

```
itd-go/
├── client.go              # Главный клиент SDK
├── config.go              # Конфигурация SDK
├── types/                 # Публичные типы и интерфейсы
│   ├── interfaces.go      # API интерфейсы
│   ├── post.go           # Типы для постов
│   ├── comment.go        # Типы для комментариев
│   ├── user.go           # Типы для пользователей
│   ├── feedTab.go        # Enum для типов ленты
│   └── pin.go            # Типы для значков
├── internal/
│   ├── api/              # Реализации API
│   │   ├── posts/        # API постов
│   │   ├── comments/     # API комментариев
│   │   └── user/         # API пользователей
│   ├── auth/             # Аутентификация
│   ├── transport/        # HTTP транспорт
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

- **types/interfaces.go** — все публичные API интерфейсы
- **client.go** — точка входа в SDK
- **internal/pkg/iterator/iterator.go** — базовая реализация итератора
- **internal/transport/transport.go** — HTTP клиент
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
