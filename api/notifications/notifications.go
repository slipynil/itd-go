package notifications

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/slipynil/itd-go/internal/transport"
)

// Service предоставляет методы для работы с уведомлениями ITD API.
type Service struct {
	transport *transport.Client
}

// New создаёт новый экземпляр клиента для работы с уведомлениями.
func New(t *transport.Client) *Service {
	return &Service{transport: t}
}

// NewIterator создаёт итератор для получения уведомлений.
// Параметры:
//   - limit: количество уведомлений на страницу (рекомендуется 10-20)
//
// Возвращает NotificationIterator для постраничной загрузки уведомлений.
func (s *Service) NewIterator(limit int) NotificationIterator {
	return newNotificationIterator(s, limit)
}

// getNotifications получает список уведомлений с пагинацией.
func (s *Service) getNotifications(ctx context.Context, offset int, limit int) (*notificationResponse, error) {

	path := fmt.Sprintf("/api/notifications/?limit=%d&offset=%d", limit, offset)
	req, err := s.transport.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result notificationResponse
	if err := json.UnmarshalRead(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
