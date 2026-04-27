package types

import "time"

type Notification struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	TargetType string    `json:"targetType"`
	TargetID   string    `json:"targetId"`
	Preview    string    `json:"preview"`
	ReadAt     time.Time `json:"readAt"`
	CreatedAt  time.Time `json:"createdAt"`
	Actor      struct {
		ID           string `json:"id"`
		DisplayName  string `json:"displayName"`
		Username     string `json:"username"`
		Avatar       string `json:"avatar"`
		IsFollowing  bool   `json:"isFollowing"`
		IsFollowedBy bool   `json:"isFollowedBy"`
	} `json:"actor"`
	Read bool `json:"read"`
}
