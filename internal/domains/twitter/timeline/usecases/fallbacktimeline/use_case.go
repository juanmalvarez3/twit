package fallbacktimeline

type UseCase struct {
	followsService FollowerService
	tweetService   TweetService
	publisher      Publisher
}

func New(followsService FollowerService, tweetService TweetService, publisher Publisher) UseCase {
	return UseCase{
		followsService: followsService,
		tweetService:   tweetService, publisher: publisher}
}
