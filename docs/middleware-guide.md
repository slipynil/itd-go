# Руководство по Middleware в itd-go

## Что такое Middleware?

**Middleware** (промежуточное ПО) — это слой кода, который находится между вашим приложением и внешним сервисом (в нашем случае — API итд.com). Middleware перехватывает HTTP-запросы и ответы, позволяя добавлять дополнительную логику без изменения основного кода.

### Аналогия из жизни

Представьте, что вы отправляете письмо:

1. **Вы пишете письмо** (ваш код создаёт HTTP-запрос)
2. **Секретарь кладёт письмо в конверт** (middleware добавляет заголовки)
3. **Охранник проверяет адрес** (middleware валидирует запрос)
4. **Почтальон доставляет письмо** (базовый HTTP-транспорт отправляет запрос)
5. **Получаете ответ обратно** (HTTP-ответ проходит через те же слои в обратном порядке)

Каждый из этих шагов (2-4) — это middleware. Они выполняют свою работу, не зная о существовании друг друга.

## Как это работает в Go?

### Интерфейс http.RoundTripper

В Go все HTTP-запросы выполняются через интерфейс `http.RoundTripper`:

```go
type RoundTripper interface {
    RoundTrip(*http.Request) (*http.Response, error)
}
```

Это единственный метод, который нужно реализовать. Он принимает запрос и возвращает ответ (или ошибку).

### Паттерн "Обёртка" (Wrapper Pattern)

Каждый middleware — это обёртка вокруг другого `RoundTripper`:

```go
type MyMiddleware struct {
    base http.RoundTripper  // следующий слой в цепочке
}

func (m *MyMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    // 1. Делаем что-то ДО отправки запроса
    req.Header.Set("X-Custom", "value")
    
    // 2. Передаём запрос следующему слою
    resp, err := m.base.RoundTrip(req)
    
    // 3. Делаем что-то ПОСЛЕ получения ответа
    if err == nil {
        fmt.Println("Статус:", resp.StatusCode)
    }
    
    return resp, err
}
```

## Middleware в itd-go

### Архитектура цепочки

В нашем проекте используется три middleware, которые выстраиваются в цепочку:

```
Ваш код
   ↓
[authMiddleware]           ← Добавляет Bearer токен
   ↓
[retryMiddleware]          ← Повторяет запрос при 429
   ↓
[statusCheckMiddleware]    ← Проверяет статус код
   ↓
[http.DefaultTransport]    ← Отправляет запрос по сети
   ↓
API итд.com
```

**Важно:** Запрос проходит сверху вниз, а ответ — снизу вверх через те же слои.

### 1. authMiddleware — Аутентификация

**Файл:** `internal/transport/middleware.go:14-32`

**Что делает:** Добавляет заголовок `Authorization: Bearer <token>` к каждому запросу.

```go
type authMiddleware struct {
    base     http.RoundTripper  // следующий слой
    provider auth.Provider      // откуда брать токен
}

func (m *authMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    // Получаем access token
    token, err := m.provider.GetAccessToken(req.Context())
    if err != nil {
        return nil, err
    }

    // Клонируем запрос (важно! нельзя изменять оригинал)
    req = req.Clone(req.Context())
    req.Header.Set("Authorization", "Bearer "+token)

    // Передаём дальше по цепочке
    return m.base.RoundTrip(req)
}
```

**Зачем клонировать запрос?**  
В Go `*http.Request` может использоваться повторно (например, при retry). Изменение оригинального запроса может привести к race condition. Клонирование создаёт безопасную копию.

**Когда срабатывает:**  
На каждом запросе, самым первым (верхний слой цепочки).

### 2. retryMiddleware — Повторные попытки

**Файл:** `internal/transport/middleware.go:54-104`

**Что делает:** Автоматически повторяет запрос при получении ошибки 429 (Too Many Requests).

```go
type retryMiddleware struct {
    base       http.RoundTripper
    maxRetries int            // сколько раз повторять
    baseDelay  time.Duration  // начальная задержка
}

func (m *retryMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    var lastErr error

    for attempt := 0; attempt <= m.maxRetries; attempt++ {
        // Пытаемся выполнить запрос
        resp, err := m.base.RoundTrip(req)

        // Если успешно или это не 429 — возвращаем результат
        if err == nil {
            return resp, nil
        }
        if !errors.Is(err, itderrors.ErrRateLimited) {
            return resp, err
        }

        lastErr = err

        // Если это последняя попытка — выходим
        if attempt == m.maxRetries {
            break
        }

        // Exponential backoff: 1s → 2s → 4s → 8s
        delay := m.baseDelay * time.Duration(1<<uint(attempt))

        fmt.Printf("[itd-go] Rate limit exceeded (429), retry %d/%d after %v\n",
            attempt+1, m.maxRetries, delay)

        // Ждём с учётом контекста
        select {
        case <-time.After(delay):
            // Продолжаем
        case <-req.Context().Done():
            // Контекст отменён — прерываем
            return nil, req.Context().Err()
        }
    }

    return nil, lastErr
}
```

**Exponential backoff:**  
Задержка увеличивается экспоненциально: 1s, 2s, 4s, 8s...  
Формула: `baseDelay * 2^attempt`

**Зачем `select` с контекстом?**  
Если пользователь отменил операцию (например, Ctrl+C), мы должны немедленно прервать ожидание, а не продолжать retry.

**Когда срабатывает:**  
Только при ошибке 429. Другие ошибки (401, 404, 500) возвращаются сразу.

### 3. statusCheckMiddleware — Проверка статуса

**Файл:** `internal/transport/middleware.go:34-52`

**Что делает:** Проверяет HTTP статус код и преобразует 4xx/5xx в ошибки Go.

```go
type statusCheckMiddleware struct {
    base http.RoundTripper
}

func (m *statusCheckMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    // Выполняем запрос
    resp, err := m.base.RoundTrip(req)
    if err != nil {
        return nil, err
    }

    // Проверяем статус код
    if err := itderrors.CheckResponse(resp); err != nil {
        return nil, err
    }

    return resp, nil
}
```

**Зачем это нужно?**  
По умолчанию `http.Client` не считает 404 или 500 ошибкой — он просто возвращает `Response` с `StatusCode=404`. Этот middleware преобразует плохие статусы в Go-ошибки, чтобы их можно было обработать через `if err != nil`.

**Когда срабатывает:**  
После получения ответа от сервера, перед возвратом в ваш код.

## Как собирается цепочка?

**Файл:** `internal/transport/transport.go:86-115`

```go
func buildTransport(cfg Config) http.RoundTripper {
    // 1. Начинаем с базового транспорта (отправка по сети)
    base := cfg.HttpClient.Transport
    if base == nil {
        base = http.DefaultTransport
    }

    // 2. Оборачиваем в statusCheck
    var transport http.RoundTripper = &statusCheckMiddleware{
        base: base,
    }

    // 3. Оборачиваем в retry (если включено)
    if cfg.MaxRetries > 0 {
        transport = &retryMiddleware{
            base:       transport,
            maxRetries: cfg.MaxRetries,
            baseDelay:  cfg.RetryDelay,
        }
    }

    // 4. Оборачиваем в auth (самый внешний слой)
    if cfg.AuthClient != nil {
        transport = &authMiddleware{
            base:     transport,
            provider: cfg.AuthClient,
        }
    }

    return transport
}
```

**Порядок важен!**  
Цепочка строится "изнутри наружу":

1. `base` (http.DefaultTransport) — самый внутренний
2. `statusCheck` оборачивает `base`
3. `retry` оборачивает `statusCheck`
4. `auth` оборачивает `retry` — самый внешний

Когда вы вызываете `client.Do(req)`, запрос проходит через них в порядке: auth → retry → statusCheck → base.

## Пример выполнения запроса

Допустим, вы вызываете:

```go
resp, err := client.Posts.GetFeed(ctx, types.FeedTabMain, 20)
```

### Шаг 1: authMiddleware

```
Запрос: GET /api/posts/feed
Заголовки: (пусто)

↓ authMiddleware добавляет токен

Запрос: GET /api/posts/feed
Заголовки: Authorization: Bearer eyJhbGc...
```

### Шаг 2: retryMiddleware

```
↓ retryMiddleware передаёт запрос дальше
  (пока ошибок нет, ничего не делает)
```

### Шаг 3: statusCheckMiddleware

```
↓ statusCheckMiddleware передаёт запрос дальше
```

### Шаг 4: http.DefaultTransport

```
↓ Отправка по сети к итд.com

← Получен ответ: 200 OK
```

### Обратный путь

```
← statusCheckMiddleware проверяет статус
  200 OK — всё хорошо, пропускает дальше

← retryMiddleware получает успешный ответ
  Ошибки нет, пропускает дальше

← authMiddleware получает ответ
  Ничего не делает с ответом, возвращает его

← Ваш код получает resp
```

## Пример с ошибкой 429

Допустим, API вернул 429 (Too Many Requests):

### Первая попытка

```
auth → retry → statusCheck → base → API
                                    ← 429 Too Many Requests
                            ← CheckResponse() возвращает ErrRateLimited
                ← retry видит ErrRateLimited
                   Ждёт 1 секунду...
```

### Вторая попытка

```
                   retry → statusCheck → base → API
                                              ← 429 Too Many Requests
                                      ← ErrRateLimited
                          ← retry видит ErrRateLimited
                             Ждёт 2 секунды...
```

### Третья попытка

```
                             retry → statusCheck → base → API
                                                        ← 200 OK
                                                ← nil error
                                    ← nil error
← Ваш код получает успешный ответ
```

## Как добавить свой middleware?

### Шаг 1: Создать структуру

```go
// loggingMiddleware логирует все запросы и ответы
type loggingMiddleware struct {
    base   http.RoundTripper
    logger *log.Logger
}
```

### Шаг 2: Реализовать RoundTrip

```go
func (m *loggingMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    // Логируем запрос
    m.logger.Printf("→ %s %s", req.Method, req.URL.Path)
    start := time.Now()

    // Выполняем запрос
    resp, err := m.base.RoundTrip(req)

    // Логируем ответ
    duration := time.Since(start)
    if err != nil {
        m.logger.Printf("← ERROR after %v: %v", duration, err)
    } else {
        m.logger.Printf("← %d after %v", resp.StatusCode, duration)
    }

    return resp, err
}
```

### Шаг 3: Добавить в цепочку

В `buildTransport()`:

```go
func buildTransport(cfg Config) http.RoundTripper {
    base := cfg.HttpClient.Transport
    if base == nil {
        base = http.DefaultTransport
    }

    // Добавляем logging как самый внутренний слой
    var transport http.RoundTripper = &loggingMiddleware{
        base:   base,
        logger: cfg.Logger,
    }

    transport = &statusCheckMiddleware{
        base: transport,
    }

    // ... остальные middleware
}
```

## Частые вопросы

### Почему нельзя изменять req напрямую?

```go
// ❌ НЕПРАВИЛЬНО
func (m *authMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    req.Header.Set("Authorization", "Bearer "+token)  // изменяет оригинал!
    return m.base.RoundTrip(req)
}

// ✅ ПРАВИЛЬНО
func (m *authMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    req = req.Clone(req.Context())  // создаём копию
    req.Header.Set("Authorization", "Bearer "+token)
    return m.base.RoundTrip(req)
}
```

**Причина:** Если retryMiddleware повторит запрос, он будет использовать тот же `*http.Request`. Изменение оригинала может привести к дублированию заголовков или race condition.

### Можно ли изменить порядок middleware?

Да, но это изменит поведение. Например, если поставить `retry` перед `auth`:

```go
// Плохой порядок
auth → statusCheck → retry → base
```

Проблема: если токен истёк, `auth` вернёт ошибку, и `retry` её не увидит (он ниже по цепочке). Правильный порядок:

```go
// Хороший порядок
auth → retry → statusCheck → base
```

Теперь `retry` может повторить запрос, если `auth` или `statusCheck` вернули ошибку.

### Зачем нужен base в каждом middleware?

Это реализация паттерна "Цепочка обязанностей" (Chain of Responsibility). Каждый middleware:

1. Делает свою работу
2. Передаёт управление следующему через `base.RoundTrip()`
3. Не знает, что находится дальше по цепочке

Это позволяет легко добавлять/удалять middleware без изменения остального кода.

### Что будет, если не вызвать base.RoundTrip()?

Запрос не будет отправлен! Например:

```go
func (m *cacheMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    // Проверяем кеш
    if cached := m.cache.Get(req.URL.String()); cached != nil {
        return cached, nil  // возвращаем из кеша, не вызывая base
    }

    // Кеша нет — идём дальше по цепочке
    return m.base.RoundTrip(req)
}
```

Это легальный паттерн для кеширования или mock-тестирования.

## Преимущества паттерна Middleware

1. **Разделение ответственности:** Каждый middleware делает одну вещь
2. **Переиспользование:** Middleware можно использовать в разных проектах
3. **Тестируемость:** Каждый middleware тестируется отдельно
4. **Гибкость:** Легко добавлять/удалять функциональность
5. **Прозрачность:** Middleware не влияет на основной код API

## Примеры использования в других проектах

### Rate Limiting

```go
type rateLimitMiddleware struct {
    base    http.RoundTripper
    limiter *rate.Limiter
}

func (m *rateLimitMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    if err := m.limiter.Wait(req.Context()); err != nil {
        return nil, err
    }
    return m.base.RoundTrip(req)
}
```

### Метрики

```go
type metricsMiddleware struct {
    base    http.RoundTripper
    counter *prometheus.CounterVec
}

func (m *metricsMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    resp, err := m.base.RoundTrip(req)
    
    status := "error"
    if err == nil {
        status = fmt.Sprintf("%d", resp.StatusCode)
    }
    
    m.counter.WithLabelValues(req.Method, status).Inc()
    return resp, err
}
```

### Кеширование

```go
type cacheMiddleware struct {
    base  http.RoundTripper
    cache Cache
}

func (m *cacheMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
    if req.Method != "GET" {
        return m.base.RoundTrip(req)  // кешируем только GET
    }

    key := req.URL.String()
    if cached := m.cache.Get(key); cached != nil {
        return cached, nil
    }

    resp, err := m.base.RoundTrip(req)
    if err == nil && resp.StatusCode == 200 {
        m.cache.Set(key, resp)
    }

    return resp, err
}
```

## Заключение

Middleware — это мощный паттерн для добавления cross-cutting concerns (сквозной функциональности) в HTTP-клиент. В itd-go он используется для:

- Автоматической аутентификации (authMiddleware)
- Обработки rate limiting (retryMiddleware)
- Преобразования HTTP-ошибок (statusCheckMiddleware)

Понимание этого паттерна позволит вам легко расширять функциональность SDK без изменения основного кода API.
