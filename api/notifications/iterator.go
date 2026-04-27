package notifications

import (
	"context"

	"github.com/slipynil/itd-go/internal/iterator"
	"github.com/slipynil/itd-go/types"
)

// NotificationIterator предоставляет интерфейс для постраничной загрузки уведомлений.
type NotificationIterator interface {
	// HasMore возвращает true, если есть ещё данные для загрузки.
	HasMore() bool
	// Next загружает и возвращает следующую страницу уведомлений.
	// Параметры:
	//   - ctx: контекст для управления временем жизни запроса
	Next(ctx context.Context) ([]*types.Notification, error)
}

func newNotificationIterator(s *Service, limit int) NotificationIterator {
	fetch := func(ctx context.Context, token iterator.PageToken) ([]*types.Notification, iterator.PageToken, bool, error) {
		result, err := s.getNotifications(ctx, token.Offset, limit)
		if err != nil {
			return nil, iterator.PageToken{}, false, err
		}
		next := iterator.PageToken{Offset: token.Offset + limit}
		return result.Notifications, next, result.HasMore, nil
	}

	return iterator.New[*types.Notification](fetch, iterator.PageToken{})
}
