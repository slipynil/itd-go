package search

import "github.com/slipynil/itd-go/types"

type topHashtagsResponse struct {
	Data struct {
		Hashtag []types.Hashtag `json:"hashtags"`
	} `json:"data"`
}

type topClansResponse struct {
	Clans []types.Clans `json:"clans"`
}

type hashtagFeedResponse struct {
	Data FeedData `json:"data"`
}

type FeedData struct {
	Hashtag    types.Hashtag `json:"hashtag"`
	Posts      []*types.Post `json:"posts"`
	Pagination struct {
		Limit      int    `json:"limit"`
		NextCursor string `json:"nextCursor"`
		HasMore    bool   `json:"hasMore"`
	} `json:"pagination"`
}
