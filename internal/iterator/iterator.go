package iterator

import "context"

// PageToken универсальный токен пагинации — может быть cursor или номером страницы.
type PageToken struct {
	// Cursor - строковый курсор для cursor-based пагинации
	Cursor string

	// Page - номер страницы для page-based пагинации
	Page int
}

// FetchFunc определяет функцию для получения следующей страницы данных.
// Принимает контекст и курсор, возвращает массив элементов, следующий курсор и флаг hasMore.
type FetchFunc[T any] func(ctx context.Context, token PageToken) (items []T, next PageToken, hasMore bool, err error)

// Iterator предоставляет интерфейс для постраничной загрузки данных.
type Iterator[T any] interface {
	// HasMore возвращает true, если есть ещё данные для загрузки.
	HasMore() bool
	// Next загружает и возвращает следующую страницу данных.
	Next(ctx context.Context) ([]T, error)
}

// paginatedIterator реализует Iterator с использованием cursor-based пагинации.
type paginatedIterator[T any] struct {
	fetch   FetchFunc[T]
	token   PageToken
	hasMore bool
}

// New создаёт новый итератор с заданной функцией получения данных.
func New[T any](fetch FetchFunc[T], startToken PageToken) Iterator[T] {
	return &paginatedIterator[T]{
		fetch:   fetch,
		token:   startToken,
		hasMore: true,
	}
}

func (i *paginatedIterator[T]) HasMore() bool {
	return i.hasMore
}

func (i *paginatedIterator[T]) Next(ctx context.Context) ([]T, error) {
	items, next, hasMore, err := i.fetch(ctx, i.token)
	if err != nil {
		return nil, err
	}

	i.token = next
	i.hasMore = hasMore

	return items, nil
}
