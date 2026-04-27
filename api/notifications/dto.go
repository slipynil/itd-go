package notifications

import (
	"github.com/slipynil/itd-go/types"
)

type notificationResponse struct {
	Notifications []*types.Notification `json:"notifications"`
	HasMore       bool                  `json:"hasMore"`
}

// readbatchRequest представляет запрос для batch-пометки уведомлений как прочитанных.
type readbatchRequest struct {
	// Ids - список ID уведомлений для пометки как прочитанных
	Ids []string `json:"ids"`
}
