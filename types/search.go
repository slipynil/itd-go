package types

type Hashtag struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PostsCount int    `json:"postsCount"`
}

type Clans struct {
	Avatar      string `json:"avatar"`
	MemberCount int    `json:"memberCount"`
}
