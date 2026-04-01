package dto

import "github.com/slipynil/itd-go/types"

type FeedResponse struct {
	Data FeedData `json:"data"`
}

type FeedData struct {
	Posts      []types.Post `json:"posts"`
	Pagination Pagination   `json:"pagination"`
}

type Pagination struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"nextCursor"`
	HasMore    bool   `json:"hasMore"`
}
