package notifications

import (
	"context"

	"github.com/slipynil/itd-go/internal/iterator"
	"github.com/slipynil/itd-go/types"
)

// Iterator предоставляет интерфейс для постраничной загрузки уведомлений.
type Iterator interface {
	// HasMore возвращает true, если есть ещё данные для загрузки.
	HasMore() bool
	// Next загружает и возвращает следующую страницу уведомлений.
	// Параметры:
	//   - ctx: контекст для управления временем жизни запроса
	Next(ctx context.Context) ([]*types.Notification, error)
}

// newNotifications создаёт итератор для получения уведомлений.
// Параметры:
//   - s: сервис для работы с API уведомлений
//   - limit: количество уведомлений на страницу
func newNotifications(s *Service, limit int) Iterator {
	fetch := func(ctx context.Context, token *iterator.PageToken) ([]*types.Notification, *iterator.PageToken, bool, error) {
		offset := 0
		if token != nil {
			offset = token.Offset
		}

		result, err := s.getNotifications(ctx, offset, limit)
		if err != nil {
			return nil, nil, false, err
		}

		var next *iterator.PageToken
		if result.HasMore {
			next = &iterator.PageToken{Offset: offset + limit}
		}

		return result.Notifications, next, result.HasMore, nil
	}

	return iterator.New[*types.Notification](fetch, nil)
}
