package gettweet

const (
	target = "use_case_get_tweet"

	getTweet = "get_tweet"
)

type UseCase struct {
	twtService TweetsService
	logger     Logger
}

func NewUseCase(twtService TweetsService, logger Logger) UseCase {
	if logger == nil {
		panic("logger cannot be nil")
	}
	
	return UseCase{
		twtService: twtService,
		logger:     logger,
	}
}
