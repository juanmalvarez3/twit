package repository

import (
	"context"
	"fmt"
	"github.com/juanmalvarez3/twit/pkg/dynamodb"
	pkgLogger "github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide() *TweetRepository {
	dynamo, err := dynamodb.Provide(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	log, err := pkgLogger.ProvideError()
	if err != nil {
		fmt.Println(err)
	}
	return NewTweetRepository(dynamo, "tweets", log)
}
