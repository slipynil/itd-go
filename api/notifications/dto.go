package notifications

import (
	"github.com/slipynil/itd-go/types"
)

type notificationResponse struct {
	Notifications []*types.Notification `json:"notifications"`
	HasMore       bool                  `json:"hasMore"`
}
