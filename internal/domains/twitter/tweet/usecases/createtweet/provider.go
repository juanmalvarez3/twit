package createtweet

import (
	"fmt"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/services"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide() UseCase {
	log, err := logger.ProvideError()
	if err != nil {
		fmt.Println(err)
		return UseCase{}
	}

	return NewUseCase(
		services.Provide(),
		log,
	)
}
