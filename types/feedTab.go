package types

// FeedTab определяет тип ленты постов.
type FeedTab string

const (
	FeedTabClan      FeedTab = "clan"
	FeedTabPopular   FeedTab = "popular"
	FeedTabFollowing FeedTab = "following"
)
