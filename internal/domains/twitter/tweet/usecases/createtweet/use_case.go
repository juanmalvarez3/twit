package createtweet

const (
	target = "use_case_create_tweet"

	createTweet = "create_tweet"
)

type UseCase struct {
	twtService TweetsService
	logger     Logger
}

func NewUseCase(twtService TweetsService, logger Logger) UseCase {
	return UseCase{
		twtService: twtService,
		logger:     logger,
	}
}
