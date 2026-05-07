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
type UserQueryResult struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"displayName"`
	Avatar         string `json:"avatar"`
	Verified       bool   `json:"verified"`
	HasNuksta      bool   `json:"hasNuksta"`
	FollowersCount int    `json:"followersCount"`
}
type SearchResult struct {
	Users    []UserQueryResult `json:"users"`
	Hashtags []Hashtag         `json:"hashtags"`
}
