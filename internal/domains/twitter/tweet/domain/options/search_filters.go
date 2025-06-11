package options

type SearchFilters struct {
	UserID  *string
	TweetID *string
}

func NewSearchFilters() SearchFilters {
	return SearchFilters{}
}

func (sf SearchFilters) WithUserID(userID string) SearchFilters {
	sf.UserID = &userID
	return sf
}
func (sf SearchFilters) WithTweetID(tweetID string) SearchFilters {
	sf.TweetID = &tweetID
	return sf
}
