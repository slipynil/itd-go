package notifications

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/slipynil/itd-go/internal/transport"
	"github.com/slipynil/itd-go/types"
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

// ListUnread возвращает список всех непрочитанных уведомлений.
// Метод загружает уведомления через итератор и останавливается при первом прочитанном,
// так как API всегда возвращает уведомления в хронологическом порядке (новые первыми).
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//
// Возвращает список непрочитанных уведомлений или ошибку при проблемах с сетью/API.
func (s *Service) ListUnread(ctx context.Context) ([]*types.Notification, error) {
	iter := s.NewIterator(20)
	var result []*types.Notification
	for iter.HasMore() {
		notifications, err := iter.Next(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range notifications {
			if n.Read {
				return result, nil
			}
			result = append(result, n)
		}
	}
	return result, nil
}

// MarkAllRead помечает все непрочитанные уведомления как прочитанные.
// Метод сначала получает список всех непрочитанных уведомлений через ListUnread,
// затем отправляет их ID одним batch-запросом для пометки как прочитанных.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//
// Возвращает ошибку при проблемах с сетью/API.
func (s *Service) MarkAllRead(ctx context.Context) error {
	notifications, err := s.ListUnread(ctx)
	if err != nil {
		return err
	}
	ids := make([]string, len(notifications))
	for i, n := range notifications {
		ids[i] = n.ID
	}
	return s.MarkRead(ctx, ids...)
}

// MarkNotificationsRead помечает переданные уведомления как прочитанные.
// Метод извлекает ID из уведомлений и отправляет batch-запрос к API.
// Удобен для случаев, когда уведомления уже загружены и нужно пометить их как прочитанные.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - notifications: список уведомлений для пометки как прочитанных
//
// Возвращает ошибку при проблемах с сетью/API.
func (s *Service) MarkNotificationsRead(ctx context.Context, notifications []*types.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	ids := make([]string, len(notifications))
	for i, n := range notifications {
		ids[i] = n.ID
	}
	return s.MarkRead(ctx, ids...)
}

// MarkRead помечает указанные уведомления как прочитанные.
// Метод принимает один или несколько ID уведомлений и отправляет batch-запрос к API.
// Если список ID пустой, метод завершается без выполнения запроса.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - ids: ID уведомлений для пометки как прочитанных
//
// Возвращает ошибку при проблемах с сетью/API.
func (s *Service) MarkRead(ctx context.Context, ids ...string) error {
	if len(ids) == 0 {
		return nil
	}
	data := readbatchRequest{
		Ids: ids,
	}
	req, err := s.transport.NewRequest(ctx, "POST", "/api/notifications/read-batch", data)
	if err != nil {
		return err
	}

	resp, err := s.transport.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Stream открывает SSE соединение для получения уведомлений в реальном времени.
// Возвращает два канала: канал уведомлений и канал ошибок.
//
// Параметры:
//   - ctx: контекст для управления временем жизни стрима
//
// Возвращает:
//   - <-chan *types.StreamNotification: канал для получения уведомлений
//   - <-chan error: канал для получения ошибок (буферизован, размер 1)
//
// Использование:
//
//	stream, errs := client.Notifications.Stream(ctx)
//	for {
//	    select {
//	    case n, ok := <-stream:
//	        if !ok {
//	            return // стрим закрыт
//	        }
//	        // обработка уведомления
//	    case err := <-errs:
//	        log.Printf("Stream error: %v", err)
//	        return // при ошибке стрим автоматически закрывается
//	    }
//	}
//
// Примечание: при получении ошибки стрим автоматически закрывается.
// Автоматическое переподключение не реализовано и должно выполняться на стороне клиента.
func (s *Service) Stream(ctx context.Context) (<-chan *types.StreamNotification, <-chan error) {
	ch := make(chan *types.StreamNotification)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)

		req, err := s.transport.NewRequest(ctx, "GET", "/api/notifications/stream", nil)
		if err != nil {
			errCh <- err
			return
		}
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := s.transport.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		var dataBuf strings.Builder
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				dataBuf.WriteString(strings.TrimPrefix(line, "data: "))
			} else if line == "" && dataBuf.Len() > 0 {
				var result types.StreamNotification
				if err := json.Unmarshal([]byte(dataBuf.String()), &result, transport.DataOptions); err != nil {
					errCh <- err
					return
				}
				// Пропускаем heartbeat события (keepalive от сервера без данных)
				if result.ID == "" {
					dataBuf.Reset()
					continue
				}
				select {
				case ch <- &result:
				case <-ctx.Done():
					return
				}
				dataBuf.Reset()
			}
		}
		if err := scanner.Err(); err != nil && ctx.Err() == nil {
			errCh <- err
		}
	}()

	return ch, errCh
}
